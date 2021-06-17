package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TylerBrock/colorjson"
	"github.com/a8m/djson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/clok/sm/aws"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

// truncateString limits the length of a string while also appending an ellipses.
func truncateString(str string, num int) string {
	short := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		short = str[0:num] + "..."
	}
	return short
}

// selectSecretNameFromList is a helper method to either bypass and return the
// `secretName` passed in via CLI flag OR retrieve a list of all secrets to allow
// for a search select by the User.
func selectSecretNameFromList(c *cli.Context) (string, error) {
	secretName := c.String("secret-id")
	if secretName == "" {
		secrets, err := sm.ListSecrets()
		if err != nil {
			PrintWarn("Error retrieving list of secrets.")
			return "", err
		}

		secretNames := make([]string, 0, len(secrets))
		for _, secret := range secrets {
			secretNames = append(secretNames, aws.StringValue(secret.Name))
		}
		sort.Strings(secretNames)

		p := &survey.Select{
			Message: "Choose a Secret to view:",
			Options: secretNames,
			Default: secretNames[0],
		}
		err = survey.AskOne(p, &secretName)
		if err != nil {
			return "", err
		}

		PrintInfo(fmt.Sprintf("Retrieving: %s", secretName))
	}
	return secretName, nil
}

// promptForEdit is a helper method providing an editor interface.
func promptForEdit(secretName string, s []byte) ([]byte, error) {
	ed := ""
	prompt := &survey.Editor{
		Message:       fmt.Sprintf("Open editor to modify '%s'?", secretName),
		FileName:      "*.json",
		Default:       string(s),
		HideDefault:   true,
		AppendDefault: true,
	}
	err := survey.AskOne(prompt, &ed, nil)
	if err != nil {
		return nil, err
	}

	return []byte(ed), nil
}

// validateAndUpdateSecretValue will take the original buffer and interactively open an
// editor to allow for updates. It will check if the original and updated buffer match, if
// so it will exit gracefully. It will also verify valid JSON and prompt if an invalid input
// is provided.
func validateAndUpdateSecretValue(secretName string, orig []byte, updateTmp []byte) ([]byte, error) {
	done := false
	for !done {
		_, err := djson.Decode(updateTmp)
		if err != nil {
			PrintWarn("invalid JSON submitted.")

			ed := false
			p1 := &survey.Confirm{
				Message: "Open to edit?",
			}
			err = survey.AskOne(p1, &ed)
			if err != nil {
				return nil, err
			}
			if ed {
				updateTmp, err = promptForEdit(secretName, updateTmp)
				if err != nil {
					return nil, cli.Exit(err, 2)
				}
				if string(orig) == strings.TrimSuffix(string(updateTmp), "\n") {
					PrintInfo("Updated value matches original. Exiting.")
					return nil, cli.Exit("", 0)
				}
			} else {
				submit := false
				p2 := &survey.Confirm{
					Message: "Continue with Submit?",
				}
				err = survey.AskOne(p2, &submit)
				if err != nil {
					return nil, err
				}
				if !submit {
					PrintWarn("Exiting without submit.")
					return nil, cli.Exit("", 0)
				}
				PrintInfo("Continuing with submit.")
				done = true
			}
		} else {
			PrintInfo("JSON validated.")
			done = true
		}
	}
	return updateTmp, nil
}

// ListSecrets CLI command to list all Secrets.
func ListSecrets(c *cli.Context) error {
	secrets, err := sm.ListSecrets()
	if err != nil {
		return cli.Exit(err, 2)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Name", "Updated", "Accessed", "Description"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Name", WidthMax: 120},
		{Name: "Updated", WidthMax: 10},
		{Name: "Accessed", WidthMax: 10},
		{
			Name:     "Description",
			WidthMax: 40,
		},
	})
	t.SortBy([]table.SortBy{
		{Name: "Name", Mode: table.Asc},
	})

	for _, secret := range secrets {
		lastdt := aws.TimeValue(secret.LastAccessedDate)
		updateddt := aws.TimeValue(secret.LastChangedDate)
		t.AppendRow([]interface{}{
			aws.StringValue(secret.Name),
			fmt.Sprintf("%d-%02d-%02d", updateddt.Year(), updateddt.Month(), updateddt.Day()),
			fmt.Sprintf("%d-%02d-%02d", lastdt.Year(), lastdt.Month(), lastdt.Day()),
			truncateString(aws.StringValue(secret.Description), 40),
		})
	}

	t.Render()

	return nil
}

// ViewSecret CLI command to view/get a Secret.
func ViewSecret(c *cli.Context) error {
	secretName, err := selectSecretNameFromList(c)
	if err != nil {
		return cli.Exit(err, 2)
	}

	secret, err := sm.GetSecret(secretName)
	if err != nil {
		return cli.Exit(err, 2)
	}

	if c.Bool("binary") {
		fmt.Println(string(secret.SecretBinary))
	} else {
		result, err := djson.Decode([]byte(aws.StringValue(secret.SecretString)))
		if err != nil {
			PrintWarn("stored string value is not valid JSON.")
			fmt.Println(aws.StringValue(secret.SecretString))
		} else {
			f := colorjson.NewFormatter()
			f.Indent = 4

			s, _ := f.Marshal(result)
			fmt.Println(string(s))
		}
	}

	return nil
}

// DescribeSecret CLI command to describe a Secret.
func DescribeSecret(c *cli.Context) error {
	secretName, err := selectSecretNameFromList(c)
	if err != nil {
		return cli.Exit(err, 2)
	}

	secret, err := sm.DescribeSecret(secretName)
	if err != nil {
		return cli.Exit(err, 2)
	}

	fmt.Println(secret.String())

	return nil
}

// EditSecret CLI command to edit a Secret.
func EditSecret(c *cli.Context) error {
	secretName, err := selectSecretNameFromList(c)
	if err != nil {
		return cli.Exit(err, 2)
	}

	secret, err := sm.GetSecret(secretName)
	if err != nil {
		return cli.Exit(err, 2)
	}

	var s []byte
	if c.Bool("binary") {
		s = secret.SecretBinary
	} else {
		result, err := djson.Decode([]byte(aws.StringValue(secret.SecretString)))
		if err != nil {
			PrintWarn("stored string value is not valid JSON.")
			s = []byte(aws.StringValue(secret.SecretString))
		} else {
			s, err = json.MarshalIndent(result, "", "    ")
			if err != nil {
				return cli.Exit(err, 2)
			}
		}
	}

	var up []byte
	up, err = promptForEdit(secretName, s)
	if err != nil {
		return cli.Exit(err, 2)
	}
	if string(s) == strings.TrimSuffix(string(up), "\n") {
		PrintInfo("Updated value matches original. Exiting.")
		return nil
	}

	var final []byte
	final, err = validateAndUpdateSecretValue(secretName, s, up)
	if err != nil {
		return cli.Exit(err, 2)
	}

	var t string
	if c.Bool("binary") {
		t = "BinarySecret"
		_, err = sm.PutSecretBinary(secretName, final)
	} else {
		t = "StringSecret"
		_, err = sm.PutSecretString(secretName, string(final))
	}

	if err != nil {
		return cli.Exit(err, 2)
	}

	PrintSuccess(fmt.Sprintf("%s %s successfully updated.", secretName, t))

	return nil
}

// CreateSecret CLI command to create a new Secret.
func CreateSecret(c *cli.Context) error {
	secretName := c.String("secret-id")
	exists, err := sm.CheckIfSecretExists(secretName)
	if err != nil {
		return cli.Exit(err, 2)
	}
	if exists {
		PrintWarn(fmt.Sprintf("'%s' already exists. Please use a different name.", secretName))
		return nil
	}

	interactive := c.Bool("interactive")
	var value []byte
	if c.String("value") == "" {
		// Assume interactive mode
		interactive = true
		value = []byte("{}")
	} else {
		value = []byte(c.String("value"))
	}

	var s []byte
	if interactive {
		result, err := djson.Decode(value)
		if err != nil {
			PrintWarn("value is not valid JSON.")
			s = value
		} else {
			s, err = json.MarshalIndent(result, "", "    ")
			if err != nil {
				return cli.Exit(err, 2)
			}
		}

		var up []byte
		up, err = promptForEdit(secretName, s)
		if err != nil {
			return cli.Exit(err, 2)
		}
		s = up
	} else {
		s = value
	}

	var t string
	if c.Bool("binary") {
		t = "BinarySecret"
		_, err = sm.CreateSecretBinary(secretName, s, c.String("description"), c.String("tags"))
	} else {
		t = "StringSecret"
		_, err = sm.CreateSecretString(secretName, string(s), c.String("description"), c.String("tags"))
	}

	if err != nil {
		return cli.Exit(err, 2)
	}

	PrintSuccess(fmt.Sprintf("%s %s successfully created.", secretName, t))

	return nil
}

// PutSecret CLI command to apply a delta to a Secret.
func PutSecret(c *cli.Context) error {
	secretName := c.String("secret-id")
	exists, err := sm.CheckIfSecretExists(secretName)
	if err != nil {
		return cli.Exit(err, 2)
	}
	if !exists {
		PrintWarn(fmt.Sprintf("'%s' does not exists. Please create the secret first.", secretName))
		return nil
	}

	interactive := c.Bool("interactive")
	var value []byte
	if c.String("value") == "" {
		// Assume interactive mode
		interactive = true

		secret, err := sm.GetSecret(secretName)
		if err != nil {
			return cli.Exit(err, 2)
		}

		if c.Bool("binary") {
			value = secret.SecretBinary
		} else {
			result, err := djson.Decode([]byte(aws.StringValue(secret.SecretString)))
			if err != nil {
				PrintWarn("stored string value is not valid JSON.")
				value = []byte(aws.StringValue(secret.SecretString))
			} else {
				value, err = json.MarshalIndent(result, "", "    ")
				if err != nil {
					return cli.Exit(err, 2)
				}
			}
		}
	} else {
		value = []byte(c.String("value"))
	}

	var final []byte
	if interactive {
		var up []byte
		up, err = promptForEdit(secretName, value)
		if err != nil {
			return cli.Exit(err, 2)
		}
		if string(value) == strings.TrimSuffix(string(up), "\n") {
			PrintInfo("Updated value matches original. Exiting.")
			return nil
		}

		final, err = validateAndUpdateSecretValue(secretName, value, up)
		if err != nil {
			return cli.Exit(err, 2)
		}
	} else {
		final = value
	}

	var t string
	if c.Bool("binary") {
		t = "BinarySecret"
		_, err = sm.PutSecretBinary(secretName, final)
	} else {
		t = "StringSecret"
		_, err = sm.PutSecretString(secretName, string(final))
	}

	if err != nil {
		return cli.Exit(err, 2)
	}

	PrintSuccess(fmt.Sprintf("%s %s successfully put new version.", secretName, t))

	return nil
}

// DeleteSecret CLI command that will delete a Secret.
func DeleteSecret(c *cli.Context) error {
	secretName := c.String("secret-id")
	exists, err := sm.CheckIfSecretExists(secretName)
	if err != nil {
		return cli.Exit(err, 2)
	}
	if !exists {
		PrintWarn(fmt.Sprintf("'%s' was not found.", secretName))
		return nil
	}

	del := false
	p1 := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to permanentaly delete '%s'?", secretName),
	}
	err = survey.AskOne(p1, &del)
	if err != nil {
		return cli.Exit(err, 2)
	}

	if !del {
		PrintInfo("Exiting without delete.")
		return nil
	}

	force := c.Bool("force")
	_, err = sm.DeleteSecret(secretName, force)
	if err != nil {
		return cli.Exit(err, 2)
	}

	PrintSuccess(fmt.Sprintf("'%s' deleted. (force: %v)", secretName, force))

	return nil
}

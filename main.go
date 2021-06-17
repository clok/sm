package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/clok/cdocs"
	"github.com/clok/sm/cmd"
	"github.com/clok/sm/info"
	"github.com/urfave/cli/v2"
)

var version string

func main() {
	// Generate the install-manpage command
	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: info.AppName,
		Hidden:  true,
	})
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:     info.AppName,
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name: info.AppRepoOwner,
			},
		},
		Copyright:            "(c) 2021 Derek Smith",
		HelpName:             info.AppName,
		Usage:                "AWS Secrets Manager CLI Tool",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				// get-secret-value
				Name:    "get",
				Aliases: []string{"view"},
				Usage:   "select from list or pass in specific secret",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "secret-id",
						Aliases: []string{"s"},
						Usage:   "Specific Secret to view, will bypass select/search",
					},
					&cli.BoolFlag{
						Name:        "binary",
						Aliases:     []string{"b"},
						Usage:       "get the SecretBinary value",
						DefaultText: info.SecretBinaryHelp,
					},
				},
				Action: cmd.ViewSecret,
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "interactive edit of a secret String Value",
				// TODO: add UsageText
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "secret-id",
						Aliases: []string{"s"},
						Usage:   "Specific Secret to edit, will bypass select/search",
					},
					&cli.BoolFlag{
						Name:        "binary",
						Aliases:     []string{"b"},
						Usage:       "get the SecretBinary value",
						DefaultText: info.SecretBinaryHelp,
					},
					// TODO: add flag for passing version stage
				},
				Action: cmd.EditSecret,
			},
			{
				// create-secret
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create new secret in Secrets Manager",
				// TODO: add UsageText
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "secret-id",
						Aliases:  []string{"s"},
						Usage:    "Secret name",
						Required: true,
					},
					&cli.BoolFlag{
						Name:        "binary",
						Aliases:     []string{"b"},
						Usage:       "get the SecretBinary value",
						DefaultText: info.SecretBinaryHelp,
					},
					&cli.StringFlag{
						Name:    "value",
						Aliases: []string{"v"},
						Usage:   "Secret Value. Will store as a string, unless binary flag is set.",
					},
					&cli.BoolFlag{
						Name:    "interactive",
						Aliases: []string{"i"},
						Usage:   "Open interactive editor to create secret value. If no 'value' is provided, an editor will be opened by default.",
					},
					&cli.StringFlag{
						Name:    "description",
						Aliases: []string{"d"},
						Usage:   "Additional description text.",
					},
					&cli.StringFlag{
						Name:    "tags",
						Aliases: []string{"t"},
						Usage:   "key=value tags (CSV list)",
					},
				},
				Action: cmd.CreateSecret,
			},
			{
				// put-secret-value
				Name:  "put",
				Usage: "non-interactive update to a specific secret",
				UsageText: `
Stores a new encrypted secret value in the specified secret. To do this, the 
operation creates a new version and attaches it to the secret. The version 
can contain a new SecretString value or a new SecretBinary value.

This will put the value to AWSCURRENT and retain one previous version 
with AWSPREVIOUS.
`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "secret-id",
						Aliases:  []string{"s"},
						Usage:    "Secret name",
						Required: true,
					},
					&cli.BoolFlag{
						Name:        "binary",
						Aliases:     []string{"b"},
						Usage:       "get the SecretBinary value",
						DefaultText: info.SecretBinaryHelp,
					},
					&cli.StringFlag{
						Name:    "value",
						Aliases: []string{"v"},
						Usage:   "Secret Value. Will store as a string, unless binary flag is set.",
					},
					&cli.BoolFlag{
						Name:    "interactive",
						Aliases: []string{"i"},
						Usage:   "Override and open interactive editor to verify and modify the new secret value.",
					},
					// TODO: add flag for passing version stage
				},
				Action: cmd.PutSecret,
			},
			{
				Name:    "delete",
				Aliases: []string{"del"},
				Usage:   "delete a specific secret",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "secret-id",
						Aliases:  []string{"s"},
						Usage:    "Specific Secret to delete",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Bypass recovery window (30 days) and immediately delete Secret.",
					},
				},
				Action: cmd.DeleteSecret,
			},
			{
				// list-secrets
				Name:   "list",
				Usage:  "display table of all secrets with meta data",
				Action: cmd.ListSecrets,
			},
			{
				// describe-secret
				Name:  "describe",
				Usage: "print description of secret to `STDOUT`",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "secret-id",
						Aliases: []string{"s"},
						Usage:   "Specific Secret to describe, will bypass select/search",
					},
				},
				Action: cmd.DescribeSecret,
			},
			im,
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version info",
				Hidden:  true,
				Action: func(c *cli.Context) error {
					fmt.Printf("%s %s (%s/%s)\n", info.AppName, version, runtime.GOOS, runtime.GOARCH)
					return nil
				},
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

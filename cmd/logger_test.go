package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/logrusorgru/aurora/v3"
)

func Test_PrintWarn(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintWarn("This is a test!")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = rescueStderr

	wantMsg := fmt.Sprintln(aurora.Red("✖ This is a test!"))
	if string(out) != wantMsg {
		t.Errorf("%#v, wanted %#v", string(out), wantMsg)
	}
}

func Test_PrintSuccess(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintSuccess("This is a test!")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = rescueStderr

	wantMsg := fmt.Sprintln(aurora.Green("✔ This is a test!"))
	if string(out) != wantMsg {
		t.Errorf("%#v, wanted %#v", string(out), wantMsg)
	}
}

func Test_PrintInfo(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintInfo("This is a test!")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = rescueStderr

	wantMsg := fmt.Sprintln(aurora.Gray(14, "➜ This is a test!"))
	if string(out) != wantMsg {
		t.Errorf("%#v, wanted %#v", string(out), wantMsg)
	}
}

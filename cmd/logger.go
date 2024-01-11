package cmd

import (
	"fmt"
	"os"

	a "github.com/logrusorgru/aurora/v4"
)

func PrintWarn(s string) {
	_, _ = fmt.Fprintln(os.Stderr, a.Sprintf(a.Red("✖ %s"), s))
}

func PrintSuccess(s string) {
	_, _ = fmt.Fprintln(os.Stderr, a.Sprintf(a.Green("✔ %s"), s))
}

func PrintInfo(s string) {
	_, _ = fmt.Fprintln(os.Stderr, a.Sprintf(a.Gray(14, "➜ %s"), s))
}

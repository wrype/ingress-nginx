// refer to https://github.com/ahmetb/kubectx
package fzf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

const (
	// EnvFZFIgnore describes the environment variable to set to disable
	// interactive context selection when fzf is installed.
	envFZFIgnore = "INGRESS_IGNORE_FZF"

	fzfExeName = "fzf"
)

// isTerminal determines if given fd is a TTY.
func isTerminal(fd *os.File) bool {
	return isatty.IsTerminal(fd.Fd())
}

// fzfInstalled determines if fzf(1) is in PATH.
func fzfInstalled() bool {
	v, _ := exec.LookPath(fzfExeName)
	return v != ""
}

// IsInteractiveMode determines if we can do choosing with fzf.
func IsInteractiveMode(stdout *os.File) bool {
	v := os.Getenv(envFZFIgnore)
	return v == "" && isTerminal(stdout) && fzfInstalled()
}

func FzfRun(stderr io.Writer) error {
	cmd := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = stderr
	cmd.Stdout = &out
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", strings.Join(os.Args, " ")))
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return errors.New("you did not choose any of urls")
	}
	success(stderr, "Open url in browser: %s", successColor.Sprint(choice))
	return browser.OpenURL(choice)
}

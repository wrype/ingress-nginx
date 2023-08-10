package fzf

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	successColor = color.New(color.FgGreen)
)

// useColors returns true if colors are force-enabled,
// false if colors are disabled, or nil for default behavior
// which is determined based on factors like if stdout is tty.
func useColors() *bool {
	tr, fa := true, false
	if os.Getenv(envForceColor) != "" {
		return &tr
	} else if os.Getenv(envNoColor) != "" {
		return &fa
	}
	return nil
}

func init() {
	colors := useColors()
	if colors == nil {
		return
	}
	if *colors {
		errorColor.EnableColor()
		warningColor.EnableColor()
		successColor.EnableColor()
	} else {
		errorColor.DisableColor()
		warningColor.DisableColor()
		successColor.DisableColor()
	}
}

func warning(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(w, warningColor.Sprint("warning: ")+format+"\n", args...)
	return err
}

func success(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(w, successColor.Sprint("✔ ")+fmt.Sprintf(format+"\n", args...))
	return err
}

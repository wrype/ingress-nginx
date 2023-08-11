package fzf

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
)

func init() {
	successColor.EnableColor()
}

func success(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintln(w, fmt.Sprintf(successColor.Sprint("✔ ")+format, args...))
	return err
}

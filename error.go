package cliutil

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// HandleError either performs a no-op if err is nil or writes the error plus
// command usage to os.Stderr and exits with a non-zero status otherwise.
func HandleError(cmd *cobra.Command, err error, prefixes ...string) {
	if err == nil {
		return
	}
	WriteError(os.Stderr, err, prefixes...)
	cmd.Usage()
	os.Exit(1)
}

// WriteError formats and writes an error message to io.Writer w. All prefixes
// are prepended to the error message and separated by a colon plus space (: ).
// Two new lone characters are printed after the error message, as it is
// assumed that command usage follows the error message.
func WriteError(w io.Writer, err error, prefixes ...string) {
	var format string
	for _, prefix := range prefixes {
		format = format + prefix + ": "
	}
	format = format + "%v\n\n"
	fmt.Fprintf(w, format, err)
}

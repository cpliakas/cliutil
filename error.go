package cliutil

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// HandleError either performs a no-op if err is nil or writes the error to
// STDERR along with the command's help text and exits with a non-zero status.
func HandleError(cmd *cobra.Command, err error, prefixes ...string) {
	if err == nil {
		return
	}

	var format string
	for _, prefix := range prefixes {
		format = format + prefix + ": "
	}
	format = format + "%v\n\n"

	fmt.Fprintf(os.Stderr, format, err)
	cmd.Usage()

	os.Exit(1)
}

package cliutil

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// HandleError either performs a no-op if err is nil or writes the error to
// STDERR along with the command's help text and exits with a non-zero status.
func HandleError(cmd *cobra.Command, err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	cmd.Help()

	os.Exit(1)
}

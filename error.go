package cliutil

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// ErrMissingArgument returns an error when a required option is empty.
type ErrMissingArgument string

func (e ErrMissingArgument) Error() string {
	return fmt.Sprintf("missing required argument: %s", string(e))
}

// ErrInvalidArgument returns an error when the passed option is not valid.
type ErrInvalidArgument string

func (e ErrInvalidArgument) Error() string {
	return fmt.Sprintf("argument not valid: %s", string(e))
}

// ErrMissingOption returns an error when a required option is empty.
type ErrMissingOption string

func (e ErrMissingOption) Error() string {
	return fmt.Sprintf("missing required option: %s", string(e))
}

// ErrInvalidOption returns an error when the passed option is not valid.
type ErrInvalidOption string

func (e ErrInvalidOption) Error() string {
	return fmt.Sprintf("option not valid: %s", string(e))
}

// ErrNotFound returns an error when the passed resource is not found.
type ErrNotFound string

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", string(e))
}

//
// Although the functions below can be useful for simple commands that either
// succeed for fail without much happening in-between, using the LeveledLogger
// with LeveledLogger.FatalIfError can provide a better user experience.
//

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
// Two new line characters are printed after the error message, as it is
// assumed that command usage follows the error message.
func WriteError(w io.Writer, err error, prefixes ...string) {
	var format string
	for _, prefix := range prefixes {
		format = format + prefix + ": "
	}
	format = format + "%v\n\n"
	fmt.Fprintf(w, format, err)
}

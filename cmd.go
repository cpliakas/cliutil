package cliutil

import (
	"fmt"
	"strings"
)

// Use formats a value for cobra.Command.Use.
func Use(command string, args ...string) (use string) {
	use = command
	for _, arg := range args {
		use = fmt.Sprintf("%s [%s]", use, strings.ToUpper(arg))
	}
	return
}

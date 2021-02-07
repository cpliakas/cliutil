// +build !windows

package cliutil_test

import (
	"fmt"
	"os"
	"syscall"
)

func sendCtrlBreak() error {
	p, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		return fmt.Errorf("error finding process: %w", err)
	}

	err = p.Signal(os.Interrupt)
	if err != nil {
		return fmt.Errorf("error sending signal: %w", err)
	}

	return nil
}

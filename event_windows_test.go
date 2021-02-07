package cliutil_test

import (
	"fmt"
	"syscall"
)

func sendCtrlBreak() error {
	d, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return fmt.Errorf("error loading dll: %w", err)
	}

	p, err := d.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return fmt.Errorf("error finding proc: %w", err)
	}

	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(syscall.Getpid()))
	if r == 0 {
		return fmt.Errorf("error generating console ctrl event: %w", err)
	}

	return nil
}

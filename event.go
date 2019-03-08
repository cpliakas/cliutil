package cliutil

import (
	"os"
	"os/signal"
	"syscall"
)

// EventListener listens for SIGINT and SIGTERM signals and notifies the
// shutdown channel if it detects that either was sent.
func EventListener() <-chan bool {
	shutdown := make(chan bool)

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case <-ch:
				shutdown <- true
				break
			}
		}
	}()

	return shutdown
}

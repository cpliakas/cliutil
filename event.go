package cliutil

import (
	"os"
	"os/signal"
	"syscall"
)

// EventListener listens for SIGINT and SIGTERM signals and notifies the
// shutdown channel if it detects that either was sent.
func EventListener(sig chan os.Signal) <-chan bool {
	shutdown := make(chan bool)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sig:
				shutdown <- true
				break
			}
		}
	}()

	return shutdown
}

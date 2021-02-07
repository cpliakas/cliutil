package cliutil

import (
	"os"
	"os/signal"
	"syscall"
)

// EventListener listens for SIGINT and SIGTERM signals and notifies the
// shutdown channel if it detects that either was sent.
type EventListener struct {
	signal   chan os.Signal
	shutdown chan bool
}

// NewEventListener returns an EventListener with the channels initialized.
func NewEventListener() *EventListener {
	return &EventListener{
		signal:   make(chan os.Signal),
		shutdown: make(chan bool),
	}
}

// Run runs the event listener in a goroutine and sends e message to
// EventListener.shutdown if a SIGINT or SIGTERM signal is detected.
func (e *EventListener) Run() *EventListener {
	signal.Notify(e.signal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-e.signal:
				e.shutdown <- true
				break
			}
		}
	}()

	return e
}

// Wait waits for EventListener.shutdown to receive a message.
func (e *EventListener) Wait() {
	<-e.shutdown
}

// StopSignal stops relaying incoming signals to EventListener.signal.
func (e *EventListener) StopSignal() {
	signal.Stop(e.signal)
}

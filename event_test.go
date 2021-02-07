package cliutil_test

import (
	"testing"
	"time"

	"github.com/cpliakas/cliutil"
)

func TestHandleEventListener(t *testing.T) {
	e := cliutil.NewEventListener().Run()
	defer e.StopSignal()

	// TODO Do we need a ready in Run?
	time.Sleep(100 * time.Millisecond)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(3 * time.Second)
		timeout <- true
	}()

	done := make(chan bool, 1)
	fail := make(chan error, 1)
	go func() {
		if err := sendCtrlBreak(); err != nil {
			fail <- err
			return
		}
		e.Wait()
		done <- true
	}()

	select {
	case <-done:
		break
	case err := <-fail:
		t.Fatal(err)
	case <-timeout:
		t.Fatal("timeout waiting for signal")
	}
}

func interrupt() {

}

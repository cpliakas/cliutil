package cliutil_test

import (
	"syscall"
	"testing"

	"github.com/cpliakas/cliutil"
)

func TestHandleEventListener(t *testing.T) {
	e := cliutil.NewEventListener().Run()
	defer e.StopSignal()
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	e.Wait()
}

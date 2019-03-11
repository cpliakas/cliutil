package cliutil_test

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/cpliakas/cliutil"
)

func TestHandleEventListener(t *testing.T) {
	sig := make(chan os.Signal)

	shutdown := cliutil.EventListener(sig)
	defer signal.Stop(sig)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	<-shutdown
}

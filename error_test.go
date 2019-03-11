package cliutil_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Command used for testing",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("test") },
}

func TestHandleErrorNil(t *testing.T) {
	cliutil.HandleError(testCmd, nil)
}

// See https://talks.golang.org/2014/testing.slide#23
func TestHandleError(t *testing.T) {
	if os.Getenv("CLIUTIL_TEST_HANDLE_ERROR") == "1" {
		cliutil.HandleError(testCmd, errors.New("because reasons"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestHandleError")
	cmd.Env = append(os.Environ(), "CLIUTIL_TEST_HANDLE_ERROR=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestWriteError(t *testing.T) {
	want := "prefix 1: prefix 2: because reasons\n\n"

	buf := bytes.NewBufferString("")
	cliutil.WriteError(buf, errors.New("because reasons"), "prefix 1", "prefix 2")

	have := buf.String()
	if have != want {
		t.Errorf("have %q, want %q", have, want)
	}
}

func TestWriteErrorNoPrefix(t *testing.T) {
	want := "because reasons\n\n"

	buf := bytes.NewBufferString("")
	cliutil.WriteError(buf, errors.New("because reasons"))

	have := buf.String()
	if have != want {
		t.Errorf("have %q, want %q", have, want)
	}
}

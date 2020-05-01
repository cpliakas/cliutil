package cliutil_test

import (
	"testing"

	"github.com/cpliakas/cliutil"
)

func TestUse(t *testing.T) {
	tests := []struct {
		command string
		args    []string
		want    string
	}{
		{"command1", []string{}, "command1"},
		{"command2", []string{"arg1"}, "command2 [ARG1]"},
		{"command3", []string{"arg1", "ArG2"}, "command3 [ARG1] [ARG2]"},
	}

	for _, tt := range tests {
		have := cliutil.Use(tt.command, tt.args...)
		if have != tt.want {
			t.Errorf("have %q, want %q", have, tt.want)
		}
	}
}

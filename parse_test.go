package cliutil_test

import (
	"testing"

	"github.com/cpliakas/cliutil"
)

func TestParseKeyValue(t *testing.T) {
	s := `time="2017-05-30T19:02:08-05:00" level=info msg="some log message" no_value`
	m := cliutil.ParseKeyValue(s)

	tests := []struct {
		key  string
		want string
	}{
		{"time", "2017-05-30T19:02:08-05:00"},
		{"level", "info"},
		{"msg", "some log message"},
		{"no_value", ""},
	}

	for _, tt := range tests {
		if m[tt.key] != tt.want {
			t.Errorf("have %q, want %q", m[tt.key], tt.want)
		}
	}
}

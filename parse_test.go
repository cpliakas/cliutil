package cliutil_test

import (
	"testing"

	"github.com/cpliakas/cliutil"
	"github.com/go-test/deep"
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
			t.Errorf("got %q, want %q", m[tt.key], tt.want)
		}
	}
}

func TestParseIntSlice(t *testing.T) {
	tests := []struct {
		s     string
		ex    []int
		exErr bool
	}{
		{"1,2,  3 ", []int{1, 2, 3}, false},
		{"1,2:4 ", []int{1, 2, 3, 4}, false},
		{"3:1", []int{3, 2, 1}, false},
		{"", []int{}, false},
		{"1,2,three", []int{}, true},
		{"1:two", []int{}, true},
		{"1:2:3", []int{}, true},
	}

	for _, tt := range tests {
		actual, _ := cliutil.ParseIntSlice(tt.s)
		if diff := deep.Equal(actual, tt.ex); diff != nil {
			t.Error(diff)
		}
	}
}

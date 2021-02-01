package cliutil_test

import (
	"testing"

	"github.com/cpliakas/cliutil"
)

func TestIsLetters(t *testing.T) {
	tests := []struct {
		s  string
		ex bool
	}{
		{"abc", true},
		{"ab3", false},
		{"", false},
	}

	for _, tt := range tests {
		if actual := cliutil.IsLetters(tt.s); actual != tt.ex {
			t.Errorf("got %t, want %t", actual, tt.ex)
		}
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		s  string
		ex bool
	}{
		{"123", true},
		{"1bc", false},
		{"", false},
	}

	for _, tt := range tests {
		if actual := cliutil.IsNumber(tt.s); actual != tt.ex {
			t.Errorf("got %t, want %t", actual, tt.ex)
		}
	}
}

func TestHasSpace(t *testing.T) {
	tests := []struct {
		s  string
		ex bool
	}{
		{"has space", true},
		{"nospace", false},
		{"", false},
	}

	for _, tt := range tests {
		if actual := cliutil.HasSpace(tt.s); actual != tt.ex {
			t.Errorf("got %t, want %t", actual, tt.ex)
		}
	}
}

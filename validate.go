package cliutil

import (
	"unicode"
)

// IsLetters returns true if s only contains letters.
func IsLetters(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsNumber returns true if s only contains numbers.
func IsNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// HasSpace returns true if s has a space.
func HasSpace(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

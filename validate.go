package cliutil

import (
	"fmt"
	"unicode"

	"github.com/spf13/viper"
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

// HasRequiredOptions validates that all required options were passed.
func HasRequiredOptions(cfg *viper.Viper, opts []string) error {
	for _, opt := range opts {
		val := cfg.Get(opt)

		// Panic if option is not set, logic error.
		if val == nil {
			panic(fmt.Errorf("option not defined: %s", opt))
		}

		// If val is a string, check that it is not empty.
		if s, ok := val.(string); ok && s == "" {
			return ErrMissingOption(opt)
		}

		// If val is an int, check that it is not zero.
		if i, ok := val.(int); ok && i == 0 {
			return ErrMissingOption(opt)
		}
	}

	return nil
}

// HasRequiredOption validates thatan option was passed.
func HasRequiredOption(cfg *viper.Viper, opt string) error {
	return HasRequiredOptions(cfg, []string{opt})
}

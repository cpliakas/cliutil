package cliutil

import (
	"strings"
	"unicode"
)

// ParseKeyValue parses key value pairs.
// See https://stackoverflow.com/a/44282136
func ParseKeyValue(s string) map[string]string {

	// Split the string by spaces that aren't in quotes.
	lastQuote := rune(0)
	parts := strings.FieldsFunc(s, func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	})

	// Build and return the key/value map.
	m := make(map[string]string)
	for _, part := range parts {
		p := strings.Split(part, "=")

		// Protect against values with no "=", treat them as a key.
		if len(p) < 2 {
			p = []string{p[0], ""}
		}

		// Trim quotes at the edges.
		p[0] = strings.TrimFunc(p[0], isQuotationMark)
		p[1] = strings.TrimFunc(p[1], isQuotationMark)

		m[p[0]] = p[1]
	}

	return m
}

func isQuotationMark(c rune) bool {
	return unicode.In(c, unicode.Quotation_Mark)
}

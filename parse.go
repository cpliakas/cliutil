package cliutil

import (
	"fmt"
	"math"
	"strconv"
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

// ParseIntSlice parses a slice of integers from a string. We expect integers
// and ranges to be separated by commas, e.g., 1,2,3:5 = []int{1,2,3,4,5}.
func ParseIntSlice(s string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}

	elems := strings.Split(s, ",")
	out := []int{}

	for idx, elem := range elems {
		i, err := parseRange(elem, idx)
		if err != nil {
			return []int{}, err
		}
		out = append(out, i...)
	}

	return out, nil
}

func parseRange(elem string, idx int) (out []int, err error) {
	elems := strings.Split(elem, ":")

	switch len(elems) {
	case 1:
		out = make([]int, 1)
		out[0], err = strconv.Atoi(strings.TrimSpace(elems[0]))
		if err != nil {
			err = fmt.Errorf("value at index %v is not a number: %w", idx, err)
			return
		}
	case 2:
		r := make([]int, 2)
		for iidx, elem := range elems {
			r[iidx], err = strconv.Atoi(strings.TrimSpace(elem))
			if err != nil {
				err = fmt.Errorf("range at index %v not valid: %w", idx, err)
				return
			}
		}
		out = Sequence(r[0], r[1])
	default:
		err = fmt.Errorf("range at index %v not valid: %w", idx, err)
		return
	}

	return
}

// Sequence returns a sequencce of numbers between start and end as an []int.
func Sequence(start, end int) []int {
	a := make([]int, int(math.Abs(float64(start-end)))+1)
	for i := range a {
		if start < end {
			a[i] = start + i
		} else {
			a[i] = start - i
		}
	}
	return a
}

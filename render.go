package cliutil

import (
	"encoding/json"
	"fmt"

	jmespath "github.com/jmespath/go-jmespath"
)

// FormatJSON returns pretty-printed JSON as a string.
func FormatJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "    ")
	return string(b), err
}

// FormatJSONWithFilter applies a JMESPath filter and returns pretty-printed
// JSON as a string and panics on any marshal errors.
func FormatJSONWithFilter(v interface{}, filter string) (out string, err error) {
	if filter != "" {
		if v, err = jmespath.Search(filter, v); err != nil {
			return
		}
	}
	out, err = FormatJSON(v)
	return
}

// PrintJSON writes pretty-printed JSON to STDOUT.
func PrintJSON(v interface{}) (err error) {
	var out string
	if out, err = FormatJSON(v); err == nil {
		_, err = fmt.Println(out)
	}
	return
}

// PrintJSONWithFilter applies a JMESPath filter and writes pretty-printed JSON
// to STDOUT.
func PrintJSONWithFilter(v interface{}, filter string) error {
	s, err := FormatJSONWithFilter(v, filter)
	fmt.Println(s)
	return err
}

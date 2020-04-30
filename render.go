package cliutil

import (
	"encoding/json"
	"fmt"

	jmespath "github.com/jmespath/go-jmespath"
)

// FormatJSON returns pretty-printed JSON as a string and panics on any
// marshal errors.
func FormatJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// FormatJSONWithFilter applies a JMESPath filter and returns pretty-printed
// JSON as a string and panics on any marshal errors.
func FormatJSONWithFilter(v interface{}, filter string) (s string, err error) {
	if filter != "" {
		if v, err = jmespath.Search(filter, v); err != nil {
			return
		}
	}
	s = FormatJSON(v)
	return
}

// PrintJSON writes pretty-printed JSON to STDOUT.
func PrintJSON(v interface{}) {
	fmt.Println(FormatJSON(v))
}

// PrintJSONWithFilter applies a JMESPath filter and writes pretty-printed JSON
// to STDOUT.
func PrintJSONWithFilter(v interface{}, filter string) error {
	s, err := FormatJSONWithFilter(v, filter)
	fmt.Println(s)
	return err
}

package cliutil_test

import (
	"github.com/cpliakas/cliutil"
)

type JsonData struct {
	Data string `json:"data"`
}

func ExamplePrintJSON() {
	cliutil.PrintJSON(&JsonData{Data: "test"})
	// Output: {
	//     "data": "test"
	// }
}

func ExamplePrintJSONWithFilter() {
	cliutil.PrintJSONWithFilter(&JsonData{Data: "test"}, "data")
	// Output: "test"
}

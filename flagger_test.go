package cliutil_test

import (
	"testing"

	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Input struct {
	ValueOne   string  `cliutil:"option=value-one short=O default='some value' usage='value one usage'"`
	ValueTwo   int     `cliutil:"option=value-two default=1"`
	ValueThree bool    `cliutil:"option=value-three default=true"`
	ValueFour  float64 `cliutil:"option=value-four default=3.14"`
}

func TestSetOptions(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test set optionss",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	viper := viper.New()
	flags := cliutil.NewFlagger(cmd, viper)

	input := &Input{}
	err := flags.SetOptions(input)
	if err != nil {
		t.Fatal(err)
	}

	ex1 := "some value"
	if actual := viper.GetString("value-one"); actual != ex1 {
		t.Errorf("got %q, expected %q", actual, ex1)
	}

	ex2 := 1
	if actual := viper.GetInt("value-two"); actual != ex2 {
		t.Errorf("got %v, expected %v", actual, ex2)
	}

	ex3 := true
	if actual := viper.GetBool("value-three"); actual != ex3 {
		t.Errorf("got %t, expected %t", actual, ex3)
	}

	ex4 := 3.14
	if actual := viper.GetFloat64("value-four"); actual != ex4 {
		t.Errorf("got %v, expected %v", actual, ex4)
	}
}

func TestGetOptions(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test set optionss",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	viper := viper.New()
	flags := cliutil.NewFlagger(cmd, viper)

	input := &Input{}
	err := flags.SetOptions(input)
	if err != nil {
		t.Fatal(err)
	}

	err = cliutil.GetOptions(input, viper)
	if err != nil {
		t.Fatal(err)
	}

	ex1 := "some value"
	if actual := input.ValueOne; actual != ex1 {
		t.Errorf("got %q, expected %q", actual, ex1)
	}

	ex2 := 1
	if actual := input.ValueTwo; actual != ex2 {
		t.Errorf("got %v, expected %v", actual, ex2)
	}

	ex3 := true
	if actual := input.ValueThree; actual != ex3 {
		t.Errorf("got %t, expected %t", actual, ex3)
	}

	ex4 := 3.14
	if actual := input.ValueFour; actual != ex4 {
		t.Errorf("got %v, expected %v", actual, ex4)
	}
}

package cliutil

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TagName is the name of the tag.
const TagName = "cliutil"

// InitConfig returns a *viper.Viper with an environment variable prefix set
// so that options can be passed from environment variables.
func InitConfig(envPrefix string) (c *viper.Viper) {
	c = viper.New()
	c.SetEnvPrefix(envPrefix)
	c.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	c.AutomaticEnv()
	return
}

// AddCommand adds a comand to it's parent, initializes the configuration,
// and returns a flagger to easily add options.
func AddCommand(parentCmd, cmd *cobra.Command, envPrefix string) (*viper.Viper, *Flagger) {
	parentCmd.AddCommand(cmd)
	cfg := InitConfig(envPrefix)
	return cfg, NewFlagger(cmd, cfg)
}

// Flagger is a utility that streamlines adding flags to commands.
type Flagger struct {
	cmd *cobra.Command
	cfg *viper.Viper
}

// NewFlagger returns a new Flagger with the *cobra.Command and *viper.Viper
// set as properties.
func NewFlagger(cmd *cobra.Command, cfg *viper.Viper) *Flagger {
	return &Flagger{cmd: cmd, cfg: cfg}
}

// Bool adds a local flag that accepts a boolean.
func (f *Flagger) Bool(name, shorthand string, value bool, usage string) {
	f.cmd.Flags().BoolP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentBool adds a persistent flag that accepts a boolean.
func (f *Flagger) PersistentBool(name, shorthand string, value bool, usage string) {
	f.cmd.PersistentFlags().BoolP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}

// Float64 adds a local flag that accepts a 64-bit float.
func (f *Flagger) Float64(name, shorthand string, value float64, usage string) {
	f.cmd.Flags().Float64P(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentFloat64 adds a persistent flag that accepts a 64-bit float.
func (f *Flagger) PersistentFloat64(name, shorthand string, value float64, usage string) {
	f.cmd.PersistentFlags().Float64P(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}

// Int adds a local flag that accepts an integer.
func (f *Flagger) Int(name, shorthand string, value int, usage string) {
	f.cmd.Flags().IntP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentInt adds a persistent flag that accepts an integer.
func (f *Flagger) PersistentInt(name, shorthand string, value int, usage string) {
	f.cmd.PersistentFlags().IntP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}

// String adds a local flag that accepts an string.
func (f *Flagger) String(name, shorthand, value, usage string) {
	f.cmd.Flags().StringP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentString adds a persistent flag that accepts an string.
func (f *Flagger) PersistentString(name, shorthand, value, usage string) {
	f.cmd.PersistentFlags().StringP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}

//
// Helper commands that set a value only if the option was passed.
//

// SetBoolValue sets s if the name flag is passed.
func SetBoolValue(cfg *viper.Viper, name string, b *bool) {
	if cfg.IsSet(name) {
		*b = cfg.GetBool(name)
	}
}

// SetFloat64Value sets f if the name flag is passed.
func SetFloat64Value(cfg *viper.Viper, name string, f *float64) {
	if cfg.IsSet(name) {
		*f = cfg.GetFloat64(name)
	}
}

// SetIntValue sets s if the name flag is passed.
func SetIntValue(cfg *viper.Viper, name string, i *int) {
	if cfg.IsSet(name) {
		*i = cfg.GetInt(name)
	}
}

// SetStringValue sets s if the name flag is passed.
func SetStringValue(cfg *viper.Viper, name string, s *string) {
	if cfg.IsSet(name) {
		*s = cfg.GetString(name)
	}
}

//
// Automatically set and get flags based on the cliutil tag.
//

// SetOptions sets flags based on the cliutil tag.
func (f *Flagger) SetOptions(v interface{}) error {

	// Get the reflection value and type.
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	// Resolve pointers.
	if rv.Kind() == reflect.Ptr {
		if rv.IsZero() {
			return errors.New("value is a zero value for its type")
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	// Structs only!
	if rv.Kind() != reflect.Struct {
		return errors.New("value must be a struct")
	}

	// Iterate over the struct's field.
	for idx := 0; idx < rt.NumField(); idx++ {
		rvf := rv.Field(idx)
		rtf := rt.Field(idx)

		// Skip fields that cannot interface, e.g., unexported field.
		if !rvf.CanInterface() {
			continue
		}

		// Resolve pointers.
		if rvf.Kind() == reflect.Ptr {
			if rvf.IsZero() {
				continue
			}
			rvf = rvf.Elem()
		}

		// Recurse into structs.
		if rvf.Kind() == reflect.Struct {
			err := f.SetOptions(rvf.Interface())
			if err != nil {
				return err
			}
			continue
		}

		// // Skip fields that cannot be set.
		// field := rv.FieldByName(rtf.Name)
		// if !field.CanSet() {
		// 	continue
		// }

		// Get the cliutil tag.
		tag := rtf.Tag.Get(TagName)
		if tag == "" {
			continue
		}

		// Parse the option from the tag.
		vals := ParseKeyValue(tag)
		name, ok := vals["option"]
		if !ok {
			continue
		}

		// Set the options.
		switch rvf.Kind() {
		case reflect.String:
			f.String(name, vals["short"], vals["default"], vals["usage"])
		case reflect.Int:
			var i int
			if s, ok := vals["default"]; ok {
				var err error
				i, err = strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("expecting %s", reflect.Int)
				}
			}
			f.Int(name, vals["short"], i, vals["usage"])
		case reflect.Bool:
			var b bool
			if s, ok := vals["default"]; ok {
				var err error
				b, err = strconv.ParseBool(s)
				if err != nil {
					return fmt.Errorf("expecting %s", reflect.Bool)
				}
			}
			f.Bool(name, vals["short"], b, vals["usage"])
		case reflect.Float64:
			var f64 float64
			if s, ok := vals["default"]; ok {
				var err error
				f64, err = strconv.ParseFloat(s, 64)
				if err != nil {
					return fmt.Errorf("expecting %s", reflect.Float32)
				}
			}
			f.Float64(name, vals["short"], f64, vals["usage"])
		default:
			return fmt.Errorf("type not supported: %s", rvf.Kind())
		}
	}

	return nil
}

// GetOptions gets values from cfg and sets them in item.
func GetOptions(v interface{}, cfg *viper.Viper) (err error) {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	return getOptions(rv, rt, cfg)
}

func getOptions(rv reflect.Value, rt reflect.Type, cfg *viper.Viper) error {

	// Resolve pointers.
	if rv.Kind() == reflect.Ptr {
		if rv.IsZero() {
			return errors.New("value is a zero value for its type")
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	// Structs only!
	if rv.Kind() != reflect.Struct {
		return errors.New("value must be a struct")
	}

	// Iterate over the struct's field.
	for idx := 0; idx < rt.NumField(); idx++ {
		rvf := rv.Field(idx)
		rtf := rt.Field(idx)

		// Skip fields that cannot interface, e.g., unexported field.
		if !rvf.CanInterface() {
			continue
		}

		// Resolve pointers.
		if rvf.Kind() == reflect.Ptr {
			if rvf.IsZero() {
				continue
			}
			rvf = rvf.Elem()
		}

		// Recurse into structs.
		if rvf.Kind() == reflect.Struct {
			err := getOptions(rvf, rvf.Type(), cfg)
			if err != nil {
				return err
			}
			continue
		}

		// Skip fields that cannot be set.
		field := rv.FieldByName(rtf.Name)
		if !field.CanSet() {
			continue
		}

		// Get the cliutil tag.
		tag := rtf.Tag.Get(TagName)
		if tag == "" {
			continue
		}

		// Parse the option from the tag.
		vals := ParseKeyValue(tag)
		name, ok := vals["option"]
		if !ok {
			continue
		}

		switch rvf.Kind() {
		case reflect.String:
			field.SetString(cfg.GetString(name))
		case reflect.Int:
			field.SetInt(int64(cfg.GetInt(name)))
		case reflect.Bool:
			field.SetBool(cfg.GetBool(name))
		case reflect.Float64:
			field.SetFloat(cfg.GetFloat64(name))
		default:
			return fmt.Errorf("type not supported: %s", rvf.Kind())
		}
	}

	return nil
}

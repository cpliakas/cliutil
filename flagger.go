package cliutil

import (
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
func AddCommand(parentCmd, cmd *cobra.Command, cfg *viper.Viper, envPrefix string) (*viper.Viper, *Flagger) {
	parentCmd.AddCommand(cmd)
	cfg = InitConfig(envPrefix)
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
func (f *Flagger) SetOptions(item interface{}) (err error) {
	val := reflect.ValueOf(item).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		// Recurse if we encounter an embedded type.
		// if field.Anonymous {
		// 	if err = f.SetOptions(val.Field(i).Interface()); err != nil {
		// 		return
		// 	}
		// 	continue
		// }

		// Get the cliutil tag, continue if empty.
		tag := field.Tag.Get(TagName)
		if tag == "" {
			continue
		}

		// Parse the values, which is in the format ParseKeyValue expects.
		vals := ParseKeyValue(tag)
		if name, ok := vals["option"]; ok {
			t := field.Type.String()
			switch t {
			case "string":
				f.String(name, vals["short"], vals["default"], vals["usage"])
			case "int":
				var i int
				if s, ok := vals["default"]; ok {
					i, err = strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("value not an integer: field=%q type=%q value=%q", field.Name, t, s)
					}
				}
				f.Int(name, vals["short"], i, vals["usage"])
			case "bool":
				var b bool
				if s, ok := vals["default"]; ok {
					b, err = strconv.ParseBool(s)
					if err != nil {
						return fmt.Errorf("value not a bool: field=%q type=%q value=%q", field.Name, t, s)
					}
				}
				f.Bool(name, vals["short"], b, vals["usage"])
			case "float64":
				var f64 float64
				if s, ok := vals["default"]; ok {
					f64, err = strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("value not a float64: field=%q type=%q value=%q", field.Name, t, s)
					}
				}
				f.Float64(name, vals["short"], f64, vals["usage"])
			default:
				return fmt.Errorf("type not supported: field=%q type=%q", field.Name, t)
			}
		}
	}

	return
}

// GetOptions gets values from cfg and sets them in item.
func GetOptions(item interface{}, cfg *viper.Viper) (err error) {
	val := reflect.ValueOf(item).Elem()
	ival := reflect.Indirect(reflect.ValueOf(item))

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		// Recurse if we encounter an embedded type.
		// if field.Anonymous {
		// 	if err = GetOptions(val.Field(i).Interface(), cfg); err != nil {
		// 		return
		// 	}
		// 	continue
		// }

		// Get the cliutil tag, continue if empty.
		tag := field.Tag.Get(TagName)
		if tag == "" {
			continue
		}

		// Get the field name.
		fname := ival.Type().Field(i).Name

		// Parse the values, which is in the format ParseKeyValue expects.
		vals := ParseKeyValue(tag)
		if name, ok := vals["option"]; ok {
			t := field.Type.String()
			switch t {
			case "string":
				val.FieldByName(fname).SetString(cfg.GetString(name))
			case "int":
				val.FieldByName(fname).SetInt(int64(cfg.GetInt(name)))
			case "bool":
				val.FieldByName(fname).SetBool(cfg.GetBool(name))
			case "float64":
				val.FieldByName(fname).SetFloat(cfg.GetFloat64(name))
			default:
				return fmt.Errorf("type not supported: field=%q type=%q", field.Name, t)
			}
		}
	}

	return
}

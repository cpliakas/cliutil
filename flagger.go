package cliutil

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

// IntSlice adds a local flag that accepts an integer slice.
func (f *Flagger) IntSlice(name, shorthand string, value []int, usage string) {
	f.cmd.Flags().IntSliceP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentIntSlice adds a persistent flag that accepts an integer slice.
func (f *Flagger) PersistentIntSlice(name, shorthand string, value []int, usage string) {
	f.cmd.PersistentFlags().IntSliceP(name, shorthand, value, usage)
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
func SetBoolValue(cfg *viper.Viper, name string, v *bool) {
	if cfg.IsSet(name) {
		*v = cfg.GetBool(name)
	}
}

// SetFloat64Value sets v if the name flag is passed.
func SetFloat64Value(cfg *viper.Viper, name string, v *float64) {
	if cfg.IsSet(name) {
		*v = cfg.GetFloat64(name)
	}
}

// SetIntValue sets v if the name flag is passed.
func SetIntValue(cfg *viper.Viper, name string, v *int) {
	if cfg.IsSet(name) {
		*v = cfg.GetInt(name)
	}
}

// SetIntSliceValue sets v if the name flag is passed.
func SetIntSliceValue(cfg *viper.Viper, name string, v *[]int) {
	if cfg.IsSet(name) {
		*v = cfg.GetIntSlice(name)
	}
}

// SetStringValue sets v if the name flag is passed.
func SetStringValue(cfg *viper.Viper, name string, v *string) {
	if cfg.IsSet(name) {
		*v = cfg.GetString(name)
	}
}

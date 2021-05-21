package cliutil

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

// TagName is the name of the tag.
const TagName = "cliutil"

// Err* variables contain comon errors.
var (
	ErrStructRequired    = errors.New("value must be a struct")
	ErrTypeNotSupported  = errors.New("type not supported")
	ErrFuncNotRegistered = errors.New("option type func not registered")
	ErrZeroValue         = errors.New("value is a zero value for its type")
)

var optmeta map[string]map[string]string

// SetOptionMetadata sets metadata for an option.
//
// Valid keys for meta are:
// - short
// - default
// - usage
// - func
func SetOptionMetadata(name string, meta map[string]string) {
	optmeta[name] = meta
}

func mergeMetadata(tag map[string]string) (map[string]string, error) {
	name, ok := tag["option"]
	if !ok {
		return tag, errors.New("option key required in cliutil tag")
	}

	meta, ok := optmeta[name]
	if !ok {
		return tag, nil
	}

	for k, v := range tag {
		meta[k] = v
	}

	return meta, nil
}

// OptionType is implemented by structs that set and read options.
type OptionType interface {

	// Set sets an option as a flag.
	Set(*Flagger) error

	// Read reads an option reflect.Value.
	Read(*viper.Viper, reflect.Value) error
}

// OptionTypeFunc is a definition for functions that return an OptionType.
type OptionTypeFunc func(map[string]string) OptionType

var optfn map[string]OptionTypeFunc

// RegisterOptionTypeFunc registers an OptionReadSetter by naame.
func RegisterOptionTypeFunc(name string, fn OptionTypeFunc) { optfn[name] = fn }

func init() {
	optmeta = make(map[string]map[string]string)

	optfn = map[string]OptionTypeFunc{
		"string":            NewStringOption,
		"int":               NewIntOption,
		"bool":              NewBoolOption,
		"float64":           NewFloat64Option,
		"[]int":             NewIntSliceOption,
		"map[string]string": NewKeyValueOption,
		"boolstring":        NewBoolStringOption,
		"ioreader":          NewIOReaderOption,
		"stdin":             NewStdinOption,
	}
}

// newOptionType returns an OptionType from the OptionTypeFunc registered to name.
func newOptionType(tag map[string]string, i interface{}) (OptionType, error) {
	var fn OptionTypeFunc

	var err error
	tag, err = mergeMetadata(tag)
	if err != nil {
		return nil, err
	}

	if name, ok := tag["func"]; ok {
		if fn, ok = optfn[name]; !ok {
			return nil, ErrFuncNotRegistered
		}
	} else {
		switch i.(type) {
		case string:
			fn = optfn["string"]
		case int:
			fn = optfn["int"]
		case bool:
			fn = optfn["bool"]
		case float64:
			fn = optfn["float64"]
		case []int:
			fn = optfn["[]int"]
		case map[string]string:
			fn = optfn["map[string]string"]
		default:
			return nil, ErrTypeNotSupported
		}
	}
	return fn(tag), nil
}

// StringOption implements Option for string options.
type StringOption struct {
	tag map[string]string
}

// NewStringOption is an OptionTypeFunc that returns a *StringOption.
func NewStringOption(tag map[string]string) OptionType { return &StringOption{tag} }

// Set implements OptionType.Set.
func (opt *StringOption) Set(f *Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements OptionType.Read.
func (opt *StringOption) Read(cfg *viper.Viper, field reflect.Value) error {
	field.SetString(cfg.GetString(opt.tag["option"]))
	return nil
}

// IntOption implements Option for int options.
type IntOption struct {
	tag map[string]string
}

// NewIntOption is an OptionTypeFunc that returns an *IntOption.
func NewIntOption(tag map[string]string) OptionType { return &IntOption{tag} }

// Set implements OptionType.Set.
func (opt *IntOption) Set(f *Flagger) (err error) {
	var v int
	if s, ok := opt.tag["default"]; ok {
		if v, err = strconv.Atoi(s); err != nil {
			return
		}
	}
	f.Int(opt.tag["option"], opt.tag["short"], v, opt.tag["usage"])
	return
}

// Read implements OptionType.Read.
func (opt *IntOption) Read(cfg *viper.Viper, field reflect.Value) error {
	field.SetInt(int64(cfg.GetInt(opt.tag["option"])))
	return nil
}

// BoolOption implements Option for bool options.
type BoolOption struct {
	tag map[string]string
}

// NewBoolOption is an OptionTypeFunc that returns a *BoolOption.
func NewBoolOption(tag map[string]string) OptionType { return &BoolOption{tag} }

// Set implements OptionType.Set.
func (opt *BoolOption) Set(f *Flagger) (err error) {
	var v bool
	if s, ok := opt.tag["default"]; ok {
		if v, err = strconv.ParseBool(s); err != nil {
			return
		}
	}
	f.Bool(opt.tag["option"], opt.tag["short"], v, opt.tag["usage"])
	return
}

// Read implements OptionType.Read.
func (opt *BoolOption) Read(cfg *viper.Viper, field reflect.Value) error {
	field.SetBool(cfg.GetBool(opt.tag["option"]))
	return nil
}

// Float64Option implements Option for float64 options.
type Float64Option struct {
	tag map[string]string
}

// NewFloat64Option is an OptionTypeFunc that returns a *Float64Option.
func NewFloat64Option(tag map[string]string) OptionType { return &Float64Option{tag} }

// Set implements OptionType.Set.
func (opt *Float64Option) Set(f *Flagger) (err error) {
	var v float64
	if s, ok := opt.tag["default"]; ok {
		if v, err = strconv.ParseFloat(s, 64); err != nil {
			return
		}
	}
	f.Float64(opt.tag["option"], opt.tag["short"], v, opt.tag["usage"])
	return
}

// Read implements OptionType.Read.
func (opt *Float64Option) Read(cfg *viper.Viper, field reflect.Value) error {
	field.SetFloat(cfg.GetFloat64(opt.tag["option"]))
	return nil
}

// IntSliceOption implements Option for []int options.
type IntSliceOption struct {
	tag map[string]string
}

// NewIntSliceOption is an OptionTypeFunc that returns an *IntSliceOption.
func NewIntSliceOption(tag map[string]string) OptionType { return &IntSliceOption{tag} }

// Set implements OptionType.Set.
func (opt *IntSliceOption) Set(f *Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements OptionType.Read.
func (opt *IntSliceOption) Read(cfg *viper.Viper, field reflect.Value) error {
	v, err := ParseIntSlice(cfg.GetString(opt.tag["option"]))
	for _, val := range v {
		field.Set(reflect.Append(field, reflect.ValueOf(val)))
	}
	return err
}

// BoolStringOption implements Option for string options.
type BoolStringOption struct {
	tag map[string]string
}

// NewBoolStringOption is an OptionTypeFunc that returns a *BoolStringOption.
func NewBoolStringOption(tag map[string]string) OptionType { return &BoolStringOption{tag} }

// Set implements OptionType.Set.
func (opt *BoolStringOption) Set(f *Flagger) (err error) {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements OptionType.Read.
func (opt *BoolStringOption) Read(cfg *viper.Viper, field reflect.Value) (err error) {
	s := cfg.GetString(opt.tag["option"])

	var v bool
	if v, err = strconv.ParseBool(s); err != nil {
		return
	}

	field.SetBool(v)
	return
}

// KeyValueOption implements Option for []int options.
type KeyValueOption struct {
	tag map[string]string
}

// NewKeyValueOption is an OptionTypeFunc that returns an *KeyValueOption.
func NewKeyValueOption(tag map[string]string) OptionType { return &KeyValueOption{tag} }

// Set implements OptionType.Set.
func (opt *KeyValueOption) Set(f *Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements OptionType.Read.
func (opt *KeyValueOption) Read(cfg *viper.Viper, field reflect.Value) error {
	v := ParseKeyValue(cfg.GetString(opt.tag["option"]))
	field.Set(reflect.ValueOf(v))
	return nil
}

// IOReaderOption implements Option for string options read from an io.Reader.
type IOReaderOption struct {
	tag map[string]string
}

// NewIOReaderOption is a OptionTypeFunc that returns a *GroupOption.
func NewIOReaderOption(tag map[string]string) OptionType { return &IOReaderOption{tag} }

// Set implements OptionType.Set.
func (opt *IOReaderOption) Set(f *Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements OptionType.Read.
func (opt *IOReaderOption) Read(cfg *viper.Viper, field reflect.Value) error {
	if !cfg.IsSet(opt.tag["option"]) {
		return nil
	}

	u, err := url.Parse(cfg.GetString(opt.tag["option"]))
	if err != nil {
		return fmt.Errorf("error parsing uri: %w", err)
	}

	r, err := getReadCloser(u)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}

	field.SetString(string(b))
	return nil
}

func getReadCloser(u *url.URL) (io.ReadCloser, error) {
	switch u.Scheme {
	case "file", "":
		return os.Open(u.Path)
	case "c":
		return os.Open(u.String())
	case "http", "https":
		resp, err := http.Get(u.String())
		return resp.Body, err
	default:
		return nil, fmt.Errorf("%s: scheme not supported", u.Scheme)
	}
}

// StdinOption implements Option for string options read via stdin.
type StdinOption struct {
	tag map[string]string
}

// NewStdinOption is an OptionTypeFunc that returns a *StdinOption.
func NewStdinOption(tag map[string]string) OptionType { return &StdinOption{tag} }

// Set implements OptionType.Set.
func (opt *StdinOption) Set(f *Flagger) error {
	f.String(opt.tag["option"], opt.tag["short"], opt.tag["default"], opt.tag["usage"])
	return nil
}

// Read implements OptionType.Read.
func (opt *StdinOption) Read(cfg *viper.Viper, field reflect.Value) (err error) {
	s := cfg.GetString(opt.tag["option"])
	if s == "" {
		if s, err = readStdin(); err != nil {
			return err
		}
	}
	field.SetString(s)
	return nil
}

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("error getting info for stdin: %w", err)
	}

	if stat.Size() == 0 {
		return "", nil
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("error reading stdin: %w", err)
	}

	return string(b), err
}

// SetOptions sets flags based on the cliutil tag.
func (f *Flagger) SetOptions(a interface{}) error {
	rv, rt, err := resolveStruct(a)
	if err != nil {
		return err
	}

	// Iterate over the struct's field.
	for idx := 0; idx < rt.NumField(); idx++ {

		// Get the field's reflect.Value and revlect.StructField.
		rvf, rtf, skip := resolveField(rv, rt, idx)
		if skip {
			continue
		}

		i := rvf.Interface()

		// Recurse into structs.
		if rvf.Kind() == reflect.Struct {
			if err := f.SetOptions(i); err != nil {
				return err
			}
			continue
		}

		// Parse the option from the tag.
		tag := parseTag(rtf)
		if _, ok := tag["option"]; !ok {
			continue
		}

		// Set the OptionType either from tag["func"] or the type of i.
		opt, err := newOptionType(tag, i)
		if err != nil {
			return fmt.Errorf("option %s: %w", tag["option"], err)
		}

		if err := opt.Set(f); err != nil {
			return fmt.Errorf("error setting option %s: %w", tag["option"], err)
		}
	}

	return nil
}

// GetOptions reads options from cfg into a.
//
// Deprecated: since v0.2.0. Use ReadOptions instead.
func GetOptions(a interface{}, cfg *viper.Viper) error { return ReadOptions(a, cfg) }

// ReadOptions reads options from cfg into a.
func ReadOptions(a interface{}, cfg *viper.Viper) (err error) {
	rv, rt, err := resolveStruct(a)
	if err != nil {
		return err
	}
	return readOptions(rv, rt, cfg)
}

func readOptions(rv reflect.Value, rt reflect.Type, cfg *viper.Viper) error {

	// Iterate over the struct's field.
	for idx := 0; idx < rt.NumField(); idx++ {

		// Get the field's reflect.Value and revlect.StructField.
		rvf, rtf, skip := resolveField(rv, rt, idx)
		if skip {
			continue
		}

		// Recurse into structs.
		if rvf.Kind() == reflect.Struct {
			err := readOptions(rvf, rvf.Type(), cfg)
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

		// Parse the option from the tag.
		tag := parseTag(rtf)
		if _, ok := tag["option"]; !ok {
			continue
		}

		// Get the OptionType either from tag["func"] or the type of i.
		i := rvf.Interface()
		opt, err := newOptionType(tag, i)
		if err != nil {
			return fmt.Errorf("option %s: %w", tag["option"], err)
		}

		// Read the option from cfg into field.
		if err := opt.Read(cfg, field); err != nil {
			return fmt.Errorf("error reading option %s: %w", tag["option"], err)
		}
	}

	return nil
}

// Adapted from html/template/content.go, https://github.com/spf13/cast.
// Copyright 2011 The Go Authors. All rights reserved.
// resolveStruct returns the value, after dereferencing as many times as
// necessary to reach the base type (or nil).
func resolveStruct(a interface{}) (rv reflect.Value, rt reflect.Type, err error) {
	if a == nil {
		err = fmt.Errorf("nil passed: %w", ErrStructRequired)
		return
	}

	rv = reflect.ValueOf(a)
	rt = reflect.TypeOf(a)

	if rt.Kind() != reflect.Ptr {
		return
	}

	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		err = fmt.Errorf("%s passed: %w", rv.Kind(), ErrStructRequired)
		return
	}

	return
}

func resolveField(v reflect.Value, t reflect.Type, idx int) (reflect.Value, reflect.StructField, bool) {
	vf := v.Field(idx)
	tf := t.Field(idx)

	// Skip fields that cannot interface, e.g., unexported field.
	if !vf.CanInterface() {
		return vf, tf, true
	}

	// Resolve pointers, skip pointers with zero values to avoid panics.
	if vf.Kind() == reflect.Ptr {
		if vf.IsZero() {
			return vf, tf, true
		}
		vf = vf.Elem()
	}

	return vf, tf, false
}

func parseTag(f reflect.StructField) (m map[string]string) {
	m = make(map[string]string)
	tag := f.Tag.Get(TagName)
	if tag != "" {
		m = ParseKeyValue(tag)
	}
	return
}

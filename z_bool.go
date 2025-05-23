package zerocfg

import (
	"encoding/json"
	"fmt"
	"strings"
)

type boolValue bool

func newBoolValue(val bool, p *bool) Value {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strToBool(s)
	if err != nil {
		return err
	}

	*b = boolValue(v)

	return err
}

func (b *boolValue) Type() string {
	return "bool"
}

// Bool registers a boolean configuration option and returns a pointer to its value.
//
// Usage:
//
//	debug := zerocfg.Bool("debug", false, "enable debug mode")
func Bool(name string, defVal bool, desc string, opts ...OptNode) *bool {
	return Any(name, defVal, desc, newBoolValue, opts...)
}

func strToBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "", "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value %q", s)
	}
}

type boolSliceValue []bool

func newBoolSlice(val []bool, p *[]bool) Value {
	*p = val
	return (*boolSliceValue)(p)
}

func (s *boolSliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *boolSliceValue) Type() string {
	return "bools"
}

// Bools registers a slice of boolean configuration options and returns a pointer to its value.
//
// Usage:
//
//	flags := zerocfg.Bools("feature.flags", []bool{true, false}, "feature flags")
func Bools(name string, value []bool, usage string, opts ...OptNode) *[]bool {
	return Any(name, value, usage, newBoolSlice, opts...)
}

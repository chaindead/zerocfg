package zfg

import (
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

func (b *boolValue) String() string {
	if *b {
		return "true"
	}
	return "false"
}

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

package zfg

import (
	"strings"
)

type stringSliceValue []string

func newStringSlice(val []string, p *[]string) Value {
	*p = val
	return (*stringSliceValue)(p)
}

func (s *stringSliceValue) Set(val string) error {
	*s = strings.Split(val, ",")
	return nil
}

func (s *stringSliceValue) Type() string {
	return "strings"
}

func (s *stringSliceValue) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}

	return strings.Join(*s, ",")
}

func Strings(name string, value []string, usage string, opts ...OptNode) *[]string {
	return Any(name, value, usage, newStringSlice, opts...)
}

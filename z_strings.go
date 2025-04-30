package zfg

import (
	"encoding/json"
)

type stringSliceValue []string

func newStringSlice(val []string, p *[]string) Value {
	*p = val
	return (*stringSliceValue)(p)
}

func (s *stringSliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *stringSliceValue) Type() string {
	return "strings"
}

func (s *stringSliceValue) String() string {
	data, _ := json.Marshal(*s)

	return string(data)
}

func Strs(name string, defVal []string, desc string, opts ...OptNode) *[]string {
	return Any(name, defVal, desc, newStringSlice, opts...)
}

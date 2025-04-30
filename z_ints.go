package zfg

import (
	"encoding/json"
	"strconv"
	"strings"
)

type intSliceValue []int

func newIntSlice(val []int, p *[]int) Value {
	*p = val
	return (*intSliceValue)(p)
}

func (s *intSliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *intSliceValue) Type() string {
	return "ints"
}

func (s *intSliceValue) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}

	parts := make([]string, len(*s))
	for i, v := range *s {
		parts[i] = strconv.Itoa(v)
	}

	return strings.Join(parts, ",")
}

func Ints(name string, value []int, usage string, opts ...OptNode) *[]int {
	return Any(name, value, usage, newIntSlice, opts...)
}

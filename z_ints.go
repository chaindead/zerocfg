package zfg

import (
	"encoding/json"
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
	data, _ := json.Marshal(*s)

	return string(data)
}

func Ints(name string, defVal []int, desc string, opts ...OptNode) *[]int {
	return Any(name, defVal, desc, newIntSlice, opts...)
}

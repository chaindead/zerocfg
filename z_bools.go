package zfg

import (
	"encoding/json"
)

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

func (s *boolSliceValue) String() string {
	data, _ := json.Marshal(*s)

	return string(data)
}

func Bools(name string, value []bool, usage string, opts ...OptNode) *[]bool {
	return Any(name, value, usage, newBoolSlice, opts...)
}

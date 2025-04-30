package zerocfg

import (
	"encoding/json"
	"strconv"
)

type intValue int

func newIntValue(val int, p *int) Value {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func Int(name string, defVal int, desc string, opts ...OptNode) *int {
	return Any(name, defVal, desc, newIntValue, opts...)
}

type int32Value int32

func newInt32Value(val int32, p *int32) Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) Type() string {
	return "int32"
}

func (i *int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func Int32(name string, defVal int32, desc string, opts ...OptNode) *int32 {
	return Any(name, defVal, desc, newInt32Value, opts...)
}

type int64Value int64

func newInt64Value(val int64, p *int64) Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func Int64(name string, defVal int64, desc string, opts ...OptNode) *int64 {
	return Any(name, defVal, desc, newInt64Value, opts...)
}

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

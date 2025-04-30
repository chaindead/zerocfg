package zfg

import (
	"encoding/json"
	"strconv"
)

type float64Value float64

func newFloat64(val float64, p *float64) Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(val string) error {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	*f = float64Value(v)
	return nil
}

func (f *float64Value) Type() string {
	return "float64"
}

func (f *float64Value) String() string {
	return strconv.FormatFloat(float64(*f), 'g', -1, 64)
}

func Float64(name string, value float64, usage string, opts ...OptNode) *float64 {
	return Any(name, value, usage, newFloat64, opts...)
}

type float64SliceValue []float64

func newFloat64Slice(val []float64, p *[]float64) Value {
	*p = val
	return (*float64SliceValue)(p)
}

func (s *float64SliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *float64SliceValue) Type() string {
	return "float64s"
}

func (s *float64SliceValue) String() string {
	data, _ := json.Marshal(*s)
	return string(data)
}

func Float64s(name string, value []float64, usage string, opts ...OptNode) *[]float64 {
	return Any(name, value, usage, newFloat64Slice, opts...)
}

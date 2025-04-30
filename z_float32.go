package zfg

import (
	"encoding/json"
	"strconv"
)

type float32Value float32

func newFloat32(val float32, p *float32) Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(val string) error {
	v, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return err
	}
	*f = float32Value(float32(v))
	return nil
}

func (f *float32Value) Type() string {
	return "float32"
}

func (f *float32Value) String() string {
	return strconv.FormatFloat(float64(*f), 'g', -1, 32)
}

func Float32(name string, value float32, usage string, opts ...OptNode) *float32 {
	return Any(name, value, usage, newFloat32, opts...)
}

type float32SliceValue []float32

func newFloat32Slice(val []float32, p *[]float32) Value {
	*p = val
	return (*float32SliceValue)(p)
}

func (s *float32SliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *float32SliceValue) Type() string {
	return "float32s"
}

func (s *float32SliceValue) String() string {
	data, _ := json.Marshal(*s)
	return string(data)
}

func Float32s(name string, value []float32, usage string, opts ...OptNode) *[]float32 {
	return Any(name, value, usage, newFloat32Slice, opts...)
}

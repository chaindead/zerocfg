package zfg

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type anyFn[T any] func(string, T, string, ...OptNode) *T

func regSource[T any](fn anyFn[T], v T) (reg func() any, value any, source map[string]string) {
	var def T
	reg = func() any {
		return fn("name", def, "desc")
	}

	source = map[string]string{"name": ToString(v)}

	return reg, v, source
}

func Test_ValueOk(t *testing.T) {
	tests := []struct {
		name string
		init func() (ptr func() any, val any, src map[string]string)
	}{
		{
			name: "default",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int, 0)
			},
		},
		{
			name: "int",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int, 5)
			},
		},
		{
			name: "uint",
			init: func() (func() any, any, map[string]string) {
				return regSource(Uint, uint(42))
			},
		},
		{
			name: "int32",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int32, int32(123))
			},
		},
		{
			name: "uint32",
			init: func() (func() any, any, map[string]string) {
				return regSource(Uint32, uint32(456))
			},
		},
		{
			name: "int64",
			init: func() (func() any, any, map[string]string) {
				return regSource(Int64, int64(789))
			},
		},
		{
			name: "uint64",
			init: func() (func() any, any, map[string]string) {
				return regSource(Uint64, uint64(1011))
			},
		},
		{
			name: "true bool",
			init: func() (func() any, any, map[string]string) {
				return regSource(Bool, true)
			},
		},
		{
			name: "false bool",
			init: func() (func() any, any, map[string]string) {
				return regSource(Bool, false)
			},
		},
		{
			name: "bools",
			init: func() (func() any, any, map[string]string) {
				return regSource(Bools, []bool{true, false, true})
			},
		},
		{
			name: "ints",
			init: func() (func() any, any, map[string]string) {
				return regSource(Ints, []int{1, 2, 3})
			},
		},
		{
			name: "str",
			init: func() (func() any, any, map[string]string) {
				return regSource(Str, "value")
			},
		},
		{
			name: "strs",
			init: func() (func() any, any, map[string]string) {
				return regSource(Strs, []string{"a", "b", "c"})
			},
		},
		{
			name: "float32",
			init: func() (func() any, any, map[string]string) {
				return regSource(Float32, float32(3.14))
			},
		},
		{
			name: "float32s",
			init: func() (func() any, any, map[string]string) {
				return regSource(Float32s, []float32{1.1, 2.2, 3.3})
			},
		},
		{
			name: "float64",
			init: func() (func() any, any, map[string]string) {
				return regSource(Float64, 3.14159265359)
			},
		},
		{
			name: "float64s",
			init: func() (func() any, any, map[string]string) {
				return regSource(Float64s, []float64{1.1, 2.2, 3.3})
			},
		},
		{
			name: "duration",
			init: func() (func() any, any, map[string]string) {
				return regSource(Dur, 5*time.Second)
			},
		},
		{
			name: "durations",
			init: func() (func() any, any, map[string]string) {
				return regSource(Durs, []time.Duration{time.Second, 2 * time.Minute, 3 * time.Hour})
			},
		},
		{
			name: "ip",
			init: func() (func() any, any, map[string]string) {
				return regSource(ipInternal, net.ParseIP("192.168.1.1"))
			},
		},
		{
			name: "map",
			init: func() (func() any, any, map[string]string) {
				return regSource(mapInternal, map[string]any{"float": 1., "str": "val"})
			},
		},
	}

	dereference := func(v any) (any, error) {
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("not a pointer")
		}
		return val.Elem().Interface(), nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = defaultConfig()

			reg, expected, source := tt.init()
			ptr := reg()

			err := c.applyParser(source)
			require.NoError(t, err)

			actual, err := dereference(ptr)
			require.NoError(t, err)

			require.EqualValues(t, expected, actual)
		})
	}
}

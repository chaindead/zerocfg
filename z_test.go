package zfg

import (
	"fmt"
	"reflect"
	"testing"

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
				return regSource(Strings, []string{"a", "b", "c"})
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

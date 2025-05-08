package zerocfg

import (
	"net"
	neturl "net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const name = "sample.name"

type anyFn[T any] func(string, T, string, ...OptNode) *T

func regSource[T any](fn anyFn[T], v T) (reg func() any, value any, source map[string]any) {
	var def T
	reg = func() any {
		return fn(name, def, "desc")
	}

	source = map[string]any{name: v}

	return reg, v, source
}

func Test_ValueOk(t *testing.T) {
	tests := []struct {
		varType string
		init    func() (ptr func() any, val any, src map[string]any)
	}{
		{
			varType: "int",
			init: func() (func() any, any, map[string]any) {
				return regSource(Int, 5)
			},
		},
		{
			varType: "uint",
			init: func() (func() any, any, map[string]any) {
				return regSource(Uint, uint(42))
			},
		},
		{
			varType: "int32",
			init: func() (func() any, any, map[string]any) {
				return regSource(Int32, int32(123))
			},
		},
		{
			varType: "uint32",
			init: func() (func() any, any, map[string]any) {
				return regSource(Uint32, uint32(456))
			},
		},
		{
			varType: "int64",
			init: func() (func() any, any, map[string]any) {
				return regSource(Int64, int64(789))
			},
		},
		{
			varType: "uint64",
			init: func() (func() any, any, map[string]any) {
				return regSource(Uint64, uint64(1011))
			},
		},
		{
			varType: "bool true",
			init: func() (func() any, any, map[string]any) {
				return regSource(Bool, true)
			},
		},
		{
			varType: "bool false",
			init: func() (func() any, any, map[string]any) {
				return regSource(Bool, false)
			},
		},
		{
			varType: "bools",
			init: func() (func() any, any, map[string]any) {
				return regSource(Bools, []bool{true, false, true})
			},
		},
		{
			varType: "ints",
			init: func() (func() any, any, map[string]any) {
				return regSource(Ints, []int{1, 2, 3})
			},
		},
		{
			varType: "string",
			init: func() (func() any, any, map[string]any) {
				return regSource(Str, "value")
			},
		},
		{
			varType: "strings",
			init: func() (func() any, any, map[string]any) {
				return regSource(Strs, []string{"a", "b", "c"})
			},
		},
		{
			varType: "float32",
			init: func() (func() any, any, map[string]any) {
				return regSource(Float32, float32(3.14))
			},
		},
		{
			varType: "floats32",
			init: func() (func() any, any, map[string]any) {
				return regSource(Floats32, []float32{1.1, 2.2, 3.3})
			},
		},
		{
			varType: "float64",
			init: func() (func() any, any, map[string]any) {
				return regSource(Float64, 3.14159265359)
			},
		},
		{
			varType: "floats64",
			init: func() (func() any, any, map[string]any) {
				return regSource(Floats64, []float64{1.1, 2.2, 3.3})
			},
		},
		{
			varType: "duration",
			init: func() (func() any, any, map[string]any) {
				return regSource(Dur, 5*time.Second)
			},
		},
		{
			varType: "durations",
			init: func() (func() any, any, map[string]any) {
				return regSource(Durs, []time.Duration{time.Second, 2 * time.Minute, 3 * time.Hour})
			},
		},
		{
			varType: "ip",
			init: func() (func() any, any, map[string]any) {
				return regSource(ipInternal, net.ParseIP("192.168.1.1"))
			},
		},
		{
			varType: "map",
			init: func() (func() any, any, map[string]any) {
				return regSource(func(name string, value map[string]any, usage string, opts ...OptNode) *map[string]any {
					return Any(name, value, usage, newMapValue, opts...)
				}, map[string]any{"float": 1., "str": "val"})
			},
		},
		{
			varType: "url",
			init: func() (reg func() any, val any, src map[string]any) {
				defaultURLStr := "http://default.com/path"
				valStr := "https://example.com/another?query=1"

				parsedVal, err := neturl.Parse(valStr)
				require.NoError(t, err)

				reg = func() any {
					return URL(name, defaultURLStr, "url description")
				}
				return reg, parsedVal, map[string]any{name: valStr}
			},
		},
	}

	dereference := func(t *testing.T, v any) any {
		val := reflect.ValueOf(v)
		require.True(t, val.Kind() == reflect.Ptr, "val must be a pointer")
		elem := val.Elem()
		return elem.Interface()
	}

	for _, tt := range tests {
		t.Run(tt.varType, func(t *testing.T) {
			c = testConfig()

			reg, expected, source := tt.init()
			ptr := reg()

			err := Parse(newMock(source))
			require.NoError(t, err)

			actual := dereference(t, ptr)

			if tt.varType == "url" {
				actualURL, okActual := actual.(*neturl.URL)
				require.True(t, okActual, "actual should be *neturl.URL")
				expectedURL, okExpected := expected.(*neturl.URL)
				require.True(t, okExpected, "expected should be *neturl.URL")

				if expectedURL == nil {
					require.Nil(t, actualURL)
				} else {
					require.NotNil(t, actualURL)
					require.Equal(t, expectedURL.String(), actualURL.String())
				}
			} else {
				require.EqualValues(t, expected, actual)
			}

			// check Set and ToString is compatible
			node, ok := c.vs[name]
			require.True(t, ok)

			stringRepresentation := ToString(actual)
			err = node.Value.Set(stringRepresentation)
			require.NoError(t, err)

			updatedActual := dereference(t, ptr)

			if tt.varType == "url" {
				updatedActualURL, okUpdated := updatedActual.(*neturl.URL)
				require.True(t, okUpdated)
				actualURL, okActual := actual.(*neturl.URL)
				require.True(t, okActual)

				if actualURL == nil {
					require.Nil(t, updatedActualURL)
				} else {
					require.NotNil(t, updatedActualURL)
					require.Equal(t, actualURL.String(), updatedActualURL.String())
				}
			} else {
				require.Equal(t, actual, updatedActual)
			}

			// check type name
			awaitedType := strings.Split(tt.varType, " ")[0]
			require.Equal(t, awaitedType, node.Value.Type())
		})
	}
}

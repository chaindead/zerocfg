package zfg

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func conv[T, K any](s []T, fn func(T) K) []K {
	r := make([]K, len(s))
	for i, v := range s {
		r[i] = fn(v)
	}

	return r
}

func Test_ValueOk(t *testing.T) {
	const (
		name = "name_key"
		desc = "description"
		num  = 10
		str  = "string_value"
	)

	var (
		strs    = []string{str, str + "1", str + "2"}
		ints    = []int{num, num * 2, num * 3}
		intsStr = conv(ints, strconv.Itoa)
	)

	tests := []struct {
		name   string
		value  func() any
		source map[string]string
		expect any
	}{
		{
			name: "default",
			value: func() any {
				return Int(name, num, desc)
			},
			expect: num,
		},
		{
			name: "int",
			value: func() any {
				return Int(name, 0, desc)
			},
			source: map[string]string{
				name: strconv.Itoa(num),
			},
			expect: num,
		},
		{
			name: "ints",
			value: func() any {
				return Ints(name, nil, desc)
			},
			source: map[string]string{
				name: strings.Join(intsStr, ","),
			},
			expect: ints,
		},
		{
			name: "str",
			value: func() any {
				return Str(name, "", desc)
			},
			source: map[string]string{
				name: str,
			},
			expect: str,
		},
		{
			name: "strs",
			value: func() any {
				return Strings(name, nil, desc)
			},
			source: map[string]string{
				name: strings.Join(strs, ","),
			},
			expect: strs,
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
			v := tt.value()

			err := c.applyParser(tt.source)
			require.NoError(t, err)

			actual, err := dereference(v)
			require.NoError(t, err)

			require.EqualValues(t, tt.expect, actual)
		})
	}
}

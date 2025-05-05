package zerocfg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ToString returns a string representation of any value for use in configuration serialization and display.
//
// The conversion rules are as follows:
//   - If the value implements fmt.Stringer, its String() method is used.
//   - For slices or arrays whose elements implement fmt.Stringer, a JSON array of their string values is returned.
//   - For other slices, arrays, or maps, the value is marshaled to JSON.
//   - For all other types, fmt.Sprint is used.
//
// ToString is used internally by zerocfg for representing option values as strings, including when passing ToString to custom parsers and for rendering configuration output.
func ToString(v any) string {
	if nv, ok := dereference(v); ok {
		return ToString(nv)
	}

	if s, ok := stringer(v); ok {
		return s.String()
	}

	rv := reflect.ValueOf(v)
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		n := rv.Len()
		strs := make([]string, 0, n)
		for i := 0; i < n; i++ {
			s, ok := stringer(rv.Index(i).Interface())
			if !ok {
				break
			}
			strs = append(strs, s.String())
		}

		if len(strs) > 0 {
			data, _ := json.Marshal(strs)
			return string(data)
		}

		fallthrough
	case reflect.Map:
		b, _ := json.Marshal(v)
		return string(b)

	default:
		// everything else: use fmt.Sprint
		return fmt.Sprint(v)
	}
}

func stringer(v any) (fmt.Stringer, bool) {
	if s, ok := v.(fmt.Stringer); ok {
		return s, true
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr && !rv.CanAddr() {
		ptr := reflect.New(rv.Type()).Interface()
		reflect.ValueOf(ptr).Elem().Set(rv)

		if s, ok := ptr.(fmt.Stringer); ok {
			return s, true
		}
	}

	return nil, false
}

func dereference(v any) (any, bool) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return nil, false
	}

	return val.Elem().Interface(), true
}

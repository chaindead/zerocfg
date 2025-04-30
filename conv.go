package zerocfg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ToString(v any) string {
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}

	rv := reflect.ValueOf(v)
	rt := rv.Type()

	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		// check if element type implements fmt.Stringer
		stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
		if rt.Elem().Implements(stringerType) {
			n := rv.Len()
			strs := make([]string, n)
			for i := range strs {
				// safe to assert because we checked Implements
				strs[i] = rv.Index(i).Interface().(fmt.Stringer).String()
			}
			data, _ := json.Marshal(strs)
			return string(data)
		}

		// fallback: generic JSON marshal of the original slice/array
		fallthrough

	case reflect.Map:
		b, _ := json.Marshal(v)
		return string(b)

	default:
		// everything else: use fmt.Sprint
		return fmt.Sprint(v)
	}
}

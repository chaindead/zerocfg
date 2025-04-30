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

	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		data, _ := json.Marshal(v)

		return string(data)
	default:
		return fmt.Sprint(v)
	}
}

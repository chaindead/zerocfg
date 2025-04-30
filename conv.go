package zfg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ToString(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice || t.Kind() == reflect.Map {
		data, _ := json.Marshal(v)

		return string(data)
	}

	return fmt.Sprint(v)
}

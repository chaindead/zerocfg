package zerocfg

import (
	"encoding/json"
)

type mapValue map[string]any

func newMapValue(val map[string]any, p *map[string]any) Value {
	*p = val
	return (*mapValue)(p)
}

func (m *mapValue) Set(val string) error {
	for k := range *m {
		delete(*m, k)
	}

	return json.Unmarshal([]byte(val), m)
}

func (m *mapValue) Type() string {
	return "map"
}

// Map registers a map[string]any configuration option and returns the map value.
//
// Usage:
//
//	limits := zerocfg.Map("limits", map[string]any{"max": 10, "min": 1}, "map of limits")
func Map(name string, defVal map[string]any, desc string, opts ...OptNode) map[string]any {
	mptr := Any(name, defVal, desc, newMapValue, opts...)

	return *mptr
}

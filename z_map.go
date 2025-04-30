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
	for k, _ := range *m {
		delete(*m, k)
	}

	return json.Unmarshal([]byte(val), m)
}

func (m *mapValue) Type() string {
	return "map"
}

func (m *mapValue) String() string {
	if m == nil {
		return ""
	}

	data, _ := json.Marshal(*m)
	return string(data)
}

func mapInternal(name string, value map[string]any, usage string, opts ...OptNode) *map[string]any {
	return Any(name, value, usage, newMapValue, opts...)
}

func Map(name string, defVal map[string]any, desc string, opts ...OptNode) map[string]any {
	mptr := Any(name, defVal, desc, newMapValue, opts...)

	return *mptr
}

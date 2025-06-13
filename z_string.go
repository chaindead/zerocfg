package zerocfg

import "encoding/json"

type stringValue string

func newStringValue(val string, p *string) Value {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Type() string {
	return "string"
}

// Str registers a string configuration option and returns a pointer to its value.
//
// Usage:
//
//	username := zerocfg.Str("db.user", "guest", "user of database")
func Str(name string, defVal string, desc string, opts ...OptNode) *string {
	return Any(name, defVal, desc, newStringValue, opts...)
}

type stringSliceValue []string

func newStringSlice(val []string, p *[]string) Value {
	*p = val
	return (*stringSliceValue)(p)
}

func (s *stringSliceValue) Set(val string) error {
	return json.Unmarshal([]byte(val), s)
}

func (s *stringSliceValue) Type() string {
	return "strings"
}

// Strs registers a slice of string configuration options and returns a pointer to its value.
//
// Usage:
//
//	hosts := zerocfg.Strs("hosts", []string{"a", "b"}, "list of hosts")
func Strs(name string, defVal []string, desc string, opts ...OptNode) *[]string {
	return Any(name, defVal, desc, newStringSlice, opts...)
}

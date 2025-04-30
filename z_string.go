package zfg

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

func (s *stringValue) String() string {
	return string(*s)
}

func Str(name string, defVal string, desc string, opts ...OptNode) *string {
	return Any(name, defVal, desc, newStringValue, opts...)
}

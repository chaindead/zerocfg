package zfg

import "strconv"

type intValue int

func newIntValue(val int, p *int) Value {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func Int(name string, value int, usage string, opts ...OptNode) *int {
	return Any(name, value, usage, newIntValue, opts...)
}

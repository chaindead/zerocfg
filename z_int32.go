package zfg

import "strconv"

type int32Value int32

func newInt32Value(val int32, p *int32) Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*i = int32Value(v)
	return err
}

func (i *int32Value) Type() string {
	return "int32"
}

func (i *int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func Int32(name string, defVal int32, desc string, opts ...OptNode) *int32 {
	return Any(name, defVal, desc, newInt32Value, opts...)
}

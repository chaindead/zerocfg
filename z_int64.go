package zfg

import "strconv"

type int64Value int64

func newInt64Value(val int64, p *int64) Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func Int64(name string, defVal int64, desc string, opts ...OptNode) *int64 {
	return Any(name, defVal, desc, newInt64Value, opts...)
}

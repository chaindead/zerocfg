package zfg

import "strconv"

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) Value {
	*p = val
	return (*uint64Value)(p)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Type() string {
	return "uint64"
}

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func Uint64(name string, defVal uint64, desc string, opts ...OptNode) *uint64 {
	return Any(name, defVal, desc, newUint64Value, opts...)
}

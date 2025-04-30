package zfg

import "strconv"

type uint32Value uint32

func newUint32Value(val uint32, p *uint32) Value {
	*p = val
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	*i = uint32Value(v)
	return err
}

func (i *uint32Value) Type() string {
	return "uint32"
}

func (i *uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

func Uint32(name string, defVal uint32, desc string, opts ...OptNode) *uint32 {
	return Any(name, defVal, desc, newUint32Value, opts...)
}

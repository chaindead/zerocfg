package zfg

import "strconv"

type uintValue uint

func newUintValue(val uint, p *uint) Value {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

func (i *uintValue) Type() string {
	return "uint"
}

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

func Uint(name string, defVal uint, desc string, opts ...OptNode) *uint {
	return Any(name, defVal, desc, newUintValue, opts...)
}

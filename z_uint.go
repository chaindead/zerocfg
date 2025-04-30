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

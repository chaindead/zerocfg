package zerocfg

import (
	"net"
)

type ipValue net.IP

func newIPValue(val net.IP, p *net.IP) Value {
	*p = val
	return (*ipValue)(p)
}

func (ip *ipValue) Set(val string) error {
	parsed := net.ParseIP(val)
	if parsed == nil {
		return &net.ParseError{Type: "IP address", Text: val}
	}

	*ip = ipValue(parsed)
	return nil
}

func (ip *ipValue) Type() string {
	return "ip"
}

func (ip *ipValue) String() string {
	return net.IP(*ip).String()
}

func ipInternal(name string, defValue net.IP, desc string, opts ...OptNode) *net.IP {
	return Any(name, defValue, desc, newIPValue, opts...)
}

func IP(name string, defValue string, desc string, opts ...OptNode) *net.IP {
	parsed := net.ParseIP(defValue)
	if parsed == nil && defValue != "" {
		panic("bad IP address: " + defValue)
	}

	return ipInternal(name, parsed, desc, opts...)
}

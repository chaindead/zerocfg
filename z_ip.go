package zerocfg

import (
	"encoding/json"
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
	if ip == nil {
		return "<nil>"
	}
	return net.IP(*ip).String()
}

type ipSliceValue []net.IP

func newIPSlice(val []net.IP, p *[]net.IP) Value {
	*p = val
	return (*ipSliceValue)(p)
}

func (ips *ipSliceValue) Set(val string) error {
	var strs []string
	if err := json.Unmarshal([]byte(val), &strs); err != nil {
		return err
	}

	parsed := make([]net.IP, len(strs))
	for i, s := range strs {
		ip := net.ParseIP(s)
		if ip == nil {
			return &net.ParseError{Type: "IP address", Text: s}
		}
		parsed[i] = ip
	}

	*ips = ipSliceValue(parsed)
	return nil
}

func (ips *ipSliceValue) Type() string {
	return "ips"
}

func ipInternal(name string, defValue net.IP, desc string, opts ...OptNode) *net.IP {
	return Any(name, defValue, desc, newIPValue, opts...)
}

// IP registers a net.IP configuration option and returns a pointer to its value.
//
// Usage:
//
//	dbIP := zerocfg.IP("db.ip", "127.0.0.1", "database IP address")
func IP(name string, defValue string, desc string, opts ...OptNode) *net.IP {
	parsed := net.ParseIP(defValue)
	if parsed == nil && defValue != "" {
		panic("bad IP address: " + defValue)
	}

	return ipInternal(name, parsed, desc, opts...)
}

func IPs(name string, defValue []string, desc string, opts ...OptNode) *[]net.IP {
	parsed := make([]net.IP, len(defValue))
	for i, s := range defValue {
		ip := net.ParseIP(s)
		if ip == nil && s != "" {
			panic("bad IP address in slice: " + s)
		}
		parsed[i] = ip
	}

	return Any(name, parsed, desc, newIPSlice, opts...)
}

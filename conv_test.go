package zerocfg

import (
	"net"
	"testing"
	"time"
)

type customStringer struct {
	text string
}

func (c customStringer) String() string {
	return c.text
}

type ptrStringer struct {
	msg string
}

func (p *ptrStringer) String() string {
	return "ptr: " + p.msg
}

func TestToString_Nil(t *testing.T) {
	// Untyped nil
	got := ToString(nil)
	if got != "" {
		t.Errorf("ToString(nil) = %q, want empty string", got)
	}

	// Typed nil pointer
	var p *int
	gotP := ToString(p)
	if gotP != "" {
		t.Errorf("ToString((*int)(nil)) = %q, want empty string", gotP)
	}

	var pStr *string
	if got := ToString(pStr); got != "" {
		t.Errorf("ToString((*string)(nil)) = %q, want empty string", got)
	}
}

func TestToString_Dereference(t *testing.T) {
	val := 123
	p := &val
	if got := ToString(p); got != "123" {
		t.Errorf("ToString(&val) = %q, want %q", got, "123")
	}

	nested := &p
	if got := ToString(nested); got != "123" {
		t.Errorf("ToString(&nested) = %q, want %q", got, "123")
	}
}

func TestToString_Stringer(t *testing.T) {
	s := customStringer{text: "hello"}
	if got := ToString(s); got != "hello" {
		t.Errorf("ToString(customStringer) = %q, want %q", got, "hello")
	}

	ptr := &ptrStringer{msg: "world"}
	if got := ToString(ptr); got != "ptr: world" {
		t.Errorf("ToString(*ptrStringer) = %q, want %q", got, "ptr: world")
	}

	// Non-pointer value whose pointer implements Stringer
	val := ptrStringer{msg: "value"}
	if got := ToString(val); got != "ptr: value" {
		t.Errorf("ToString(ptrStringer) = %q, want %q", got, "ptr: value")
	}
}

func TestToString_SlicesAndArrays(t *testing.T) {
	// Slice of Stringers
	ipSlice := []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("10.0.0.1")}
	expected := `["192.168.1.1","10.0.0.1"]`
	if got := ToString(ipSlice); got != expected {
		t.Errorf("ToString(ipSlice) = %q, want %q", got, expected)
	}

	// Slice of primitives
	ints := []int{1, 2, 3}
	if got := ToString(ints); got != `[1,2,3]` {
		t.Errorf("ToString(ints) = %q, want %q", got, `[1,2,3]`)
	}

	// Empty slice
	empty := []string{}
	if got := ToString(empty); got != `[]` {
		t.Errorf("ToString(empty) = %q, want %q", got, `[]`)
	}
}

func TestToString_PrimitivesAndMaps(t *testing.T) {
	if got := ToString(42); got != "42" {
		t.Errorf("ToString(42) = %q, want %q", got, "42")
	}

	if got := ToString(true); got != "true" {
		t.Errorf("ToString(true) = %q, want %q", got, "true")
	}

	dur := 5 * time.Second
	if got := ToString(dur); got != "5s" {
		t.Errorf("ToString(5s) = %q, want %q", got, "5s")
	}

	m := map[string]string{"key": "val"}
	if got := ToString(m); got != `{"key":"val"}` {
		t.Errorf("ToString(map) = %q, want %q", got, `{"key":"val"}`)
	}
}

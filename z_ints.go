package zfg

import (
	"fmt"
	"strconv"
	"strings"
)

type intSliceValue []int

func newIntSlice(val []int, p *[]int) Value {
	*p = val
	return (*intSliceValue)(p)
}

func (s *intSliceValue) Set(val string) error {
	if val == "" {
		*s = nil
		return nil
	}

	parts := strings.Split(val, ",")
	slice := make([]int, 0, len(parts))

	for _, part := range parts {
		num, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return fmt.Errorf("parse int %q: %w", part, err)
		}
		slice = append(slice, int(num))
	}

	*s = slice
	return nil
}

func (s *intSliceValue) Type() string {
	return "ints"
}

func (s *intSliceValue) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}

	parts := make([]string, len(*s))
	for i, v := range *s {
		parts[i] = strconv.Itoa(v)
	}

	return strings.Join(parts, ",")
}

func Ints(name string, value []int, usage string, opts ...OptNode) *[]int {
	return Any(name, value, usage, newIntSlice, opts...)
}

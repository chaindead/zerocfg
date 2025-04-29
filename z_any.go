package zfg

func Any[T any](name string, value T, usage string, create func(T, *T) Value, opts ...OptNode) *T {
	p := new(T)
	*p = value
	c.add(name, create(value, p), usage, opts...)

	return p
}

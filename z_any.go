package zfg

func Any[T any](name string, defVal T, desc string, create func(T, *T) Value, opts ...OptNode) *T {
	p := new(T)
	*p = defVal
	c.add(name, create(defVal, p), desc, opts...)

	return p
}

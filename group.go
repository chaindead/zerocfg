package zerocfg

type Grp struct {
	prefix string
	opts   []OptNode
}

func NewGroup(prefix string, opts ...OptNode) *Grp {
	return &Grp{
		prefix: prefix,
		opts:   opts,
	}
}

func NewOptions(opts ...OptNode) *Grp {
	return &Grp{
		opts: opts,
	}
}

func (g *Grp) key(child string) string {
	if g.prefix == "" {
		return child
	}

	return g.prefix + "." + child
}

func (g *Grp) applyOpts(n *Node) {
	for _, opt := range g.opts {
		opt(n)
	}
}

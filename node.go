package zerocfg

type Node struct {
	Name        string
	Description string
	Aliases     []string
	Value       Value
	fromSource  bool
	isSecret    bool
}

type OptNode func(*Node)

func Alias(alias string) OptNode {
	return func(n *Node) {
		n.Aliases = append(n.Aliases, alias)
	}
}

func Secret() OptNode {
	return func(n *Node) {
		n.isSecret = true
	}
}

func Group(g *Grp) OptNode {
	return func(n *Node) {
		n.Name = g.key(n.Name)
		g.applyOpts(n)
	}
}

type Value interface {
	Set(string) error
	String() string
	Type() string
}

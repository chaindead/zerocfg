package zfg

type Node struct {
	Name        string
	Description string
	Aliases     []string
	Value       Value
	fromSource  bool
}

type OptNode func(*Node)

func Alias(alias string) OptNode {
	return func(n *Node) {
		n.Aliases = append(n.Aliases, alias)
	}
}

type Value interface {
	String() string
	Set(string) error
	Type() string
}

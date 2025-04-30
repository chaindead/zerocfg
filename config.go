package zerocfg

import (
	"fmt"
	"strings"

	"github.com/chaindead/zerocfg/flag"
)

type config struct {
	vs      map[string]*Node
	aliases map[string]string

	parsers []Parser
}

func defaultConfig() *config {
	return &config{
		make(map[string]*Node),
		make(map[string]string),
		[]Parser{flag.New()},
	}
}

var c = defaultConfig()

func (c *config) add(key string, v Value, usage string, opts ...OptNode) {
	n := &Node{
		Name:        key,
		Description: usage,
		Value:       v,
	}

	for _, opt := range opts {
		opt(n)
	}

	if c.vs[n.Name] != nil {
		err := fmt.Errorf("key=%q: %w", n.Name, ErrDuplicateKey)
		panic(err)
	}

	c.vs[n.Name] = n
	for _, alias := range n.Aliases {
		if c.vs[alias] != nil {
			err := fmt.Errorf("key=%q: %w", alias, ErrCollidingAlias)
			panic(err)
		}

		c.aliases[alias] = n.Name
	}
}

func (c *config) set(key string, v string) error {
	trueKey, ok := c.aliases[key]
	if ok {
		key = trueKey
	}

	n, ok := c.vs[key]
	if !ok {
		return ErrNoSuchKey
	}

	if n.fromSource {
		return nil
	}

	n.fromSource = true
	return c.vs[key].Value.Set(v)
}

func (c *config) awaited() map[string]bool {
	a := make(map[string]bool)

	for k, _ := range c.vs {
		a[k] = true
	}

	for k, _ := range c.aliases {
		a[k] = false
	}

	return a
}

func Configuration() string {
	b := strings.Builder{}

	b.WriteString("config:\n")

	for k, v := range c.vs {

		val := v.Value.String()
		if v.isSecret {
			val = "<secret>"
		}

		line := fmt.Sprintf("  -%s = %s [%s]\n", k, val, v.Description)
		b.WriteString(line)
	}

	return b.String()
}

package zfg

import (
	"errors"
	"fmt"
	"strings"
	"zfg/flag"
)

type config struct {
	vs      map[string]*Node
	aliases map[string]string

	parsers []Parser
}

var c = &config{
	make(map[string]*Node),
	make(map[string]string),
	[]Parser{flag.New()},
}

func (f *config) add(key string, v Value, usage string, opts ...OptNode) {
	n := &Node{
		Name:        key,
		Description: usage,
		Value:       v,
	}

	for _, opt := range opts {
		opt(n)
	}

	f.vs[key] = n
	for _, alias := range n.Aliases {
		f.aliases[alias] = key
	}
}

func (f *config) set(key string, v string) error {
	trueKey, ok := f.aliases[key]
	if ok {
		key = trueKey
	}

	n, ok := f.vs[key]
	if !ok {
		return errors.New("no such key")
	}

	if n.fromSource {
		return nil
	}

	n.fromSource = true
	return f.vs[key].Value.Set(v)
}

func Configuration() string {
	b := strings.Builder{}

	b.WriteString("config:\n")

	for k, v := range c.vs {
		line := fmt.Sprintf("  -%s = %s [%s]\n", k, v.Value, v.Description)
		b.WriteString(line)
	}

	return b.String()
}

package zfg

import (
	"errors"
	"strings"
	"zfg/flag"

	"github.com/goccy/go-yaml"
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
	// Create a map that will be marshaled directly to YAML
	config := make(map[string]interface{})

	// Convert configuration values to a nested map structure
	for k, v := range c.vs {
		parts := strings.Split(k, ".")
		current := config

		// Navigate through the parts to create nested maps
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]interface{})
			}
			current = current[part].(map[string]interface{})
		}

		// Set the final value
		lastPart := parts[len(parts)-1]
		current[lastPart] = v.Value.String()
	}

	// Marshal the entire structure to YAML
	out, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	return string(out)
}

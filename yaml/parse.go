package yaml

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Parser struct {
	path *string
}

func New(path *string) Parser {
	return Parser{path: path}
}

func (p Parser) Type() string {
	return "yaml"
}

func (p Parser) Parse() (map[string]string, error) {
	data, err := os.ReadFile(*p.path)
	if err != nil {
		return nil, fmt.Errorf("read yaml file: %w", err)
	}

	var settings map[string]any
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return flatten(settings), nil
}

func flatten(settings map[string]any) map[string]string {
	keys := make(map[string]string)
	flattenDFS(settings, "", keys)

	return keys
}

func flattenDFS(m map[string]any, prefix string, keys map[string]string) {
	for k, v := range m {
		newKey := k
		if prefix != "" {
			newKey = prefix + "." + k
		}

		if subMap, ok := v.(map[string]interface{}); ok {
			flattenDFS(subMap, newKey, keys)

			continue
		}

		var value string
		switch t := v.(type) {
		case []any:
			ss := make([]string, 0, len(t))
			for _, sub := range t {
				ss = append(ss, fmt.Sprintf("%v", sub))
			}

			value = strings.Join(ss, ",")
		default:
			value = fmt.Sprintf("%v", v)
		}

		keys[newKey] = value
	}
}

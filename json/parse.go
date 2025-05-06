package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Parser struct {
	path *string

	conv    func(any) string
	awaited map[string]bool
}

func New(path *string) *Parser {
	return &Parser{path: path}
}

func (p *Parser) Type() string {
	if p.path == nil || *p.path == "" {
		return "json"
	}
	return fmt.Sprintf("json[%s]", *p.path)
}

func (p *Parser) Parse(keys map[string]bool, conv func(any) string) (found, unknown map[string]string, err error) {
	if p.path == nil || *p.path == "" {
		return make(map[string]string), make(map[string]string), nil
	}

	data, err := os.ReadFile(*p.path)
	if err != nil {
		return nil, nil, fmt.Errorf("read json file %q: %w", *p.path, err)
	}

	p.conv = conv
	p.awaited = keys

	return p.parse(data)
}

func (p *Parser) parse(data []byte) (found, unknown map[string]string, err error) {
	var settings any

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err = decoder.Decode(&settings); err != nil {
		if err == io.EOF && len(data) == 0 {
			return nil, nil, fmt.Errorf("unmarshal json: empty json input")
		}
		return nil, nil, fmt.Errorf("unmarshal json: %w", err)
	}

	settingsMap, ok := settings.(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("json root is not an object")
	}

	found, unknown = p.flatten(settingsMap)
	return found, unknown, nil
}

func (p *Parser) flatten(settings map[string]any) (found, unknown map[string]string) {
	found, unknown = make(map[string]string), make(map[string]string)
	p.flattenDFS(settings, "", found, unknown)
	return found, unknown
}

func (p *Parser) flattenDFS(m map[string]any, prefix string, found, unknown map[string]string) {
	for k, v := range m {
		if v == nil {
			continue
		}

		newKey := k
		if prefix != "" {
			newKey = prefix + "." + k
		}

		if p.awaited[newKey] {
			found[newKey] = p.conv(v)
			continue
		}

		if subMap, ok := v.(map[string]any); ok {
			p.flattenDFS(subMap, newKey, found, unknown)
			continue
		}

		unknown[newKey] = p.conv(v)
	}
}

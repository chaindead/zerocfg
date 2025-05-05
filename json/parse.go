package json

import (
	"encoding/json"
	"fmt"
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
	found = make(map[string]string)
	unknown = make(map[string]string)

	if p.path == nil || *p.path == "" {
		return found, unknown, nil
	}

	data, err := os.ReadFile(*p.path)
	if err != nil {
		if os.IsNotExist(err) {
			return found, unknown, nil
		}
		return nil, nil, fmt.Errorf("read json file %q: %w", *p.path, err)
	}

	if len(data) == 0 {
		return found, unknown, fmt.Errorf("json file %q is empty", *p.path)
	}

	p.conv = conv
	p.awaited = keys

	return p.parse(data)
}

func (p *Parser) parse(data []byte) (found, unknown map[string]string, err error) {
	var settings any

	tempFile, err := os.CreateTemp("", "json-input-*.json")
	if err != nil {
		return nil, nil, fmt.Errorf("create temp file for json: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return nil, nil, fmt.Errorf("write to temp file for json: %w", err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close()
		return nil, nil, fmt.Errorf("seek temp file for json: %w", err)
	}

	decoder := json.NewDecoder(tempFile)
	decoder.UseNumber()
	decodeErr := decoder.Decode(&settings)
	closeErr := tempFile.Close()

	if decodeErr != nil {
		return nil, nil, fmt.Errorf("unmarshal json: %w", decodeErr)
	}
	if closeErr != nil {
		return nil, nil, fmt.Errorf("close temp file for json: %w", closeErr)
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
		newKey := k
		if prefix != "" {
			newKey = prefix + "." + k
		}

		_, isAwaited := p.awaited[newKey]

		if subMap, ok := v.(map[string]any); ok {
			if isAwaited {
				found[newKey] = p.conv(v)

			} else {

				isParentOfAwaited := false
				for awaitedKey := range p.awaited {
					if len(awaitedKey) > len(newKey) && awaitedKey[len(newKey)] == '.' && awaitedKey[:len(newKey)] == newKey {
						isParentOfAwaited = true
						break
					}
				}

				if !isParentOfAwaited {
					unknown[newKey] = p.conv(v)
				}

				p.flattenDFS(subMap, newKey, found, unknown)
			}
		} else if _, ok := v.([]any); ok {
			if isAwaited {
				found[newKey] = p.conv(v)
			} else {
				unknown[newKey] = p.conv(v)
			}
		} else {
			if isAwaited {
				found[newKey] = p.conv(v)
			} else {
				unknown[newKey] = p.conv(v)
			}
		}
	}
}

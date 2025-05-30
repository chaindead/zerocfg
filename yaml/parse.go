package yaml

import (
	"fmt"
	"os"

	"github.com/chaindead/zerocfg/util"
	"gopkg.in/yaml.v3"
)

type Provider struct {
	path *string

	conv    func(any) string
	awaited map[string]bool
}

func New(path *string) *Provider {
	return &Provider{path: path}
}

func (p *Provider) Type() string {
	return fmt.Sprintf("yaml[%s]", util.ShortenPath(*p.path))
}

func (p *Provider) Provide(keys map[string]bool, conv func(any) string) (found, unknown map[string]string, err error) {
	data, err := os.ReadFile(*p.path)
	if err != nil {
		return nil, nil, fmt.Errorf("read yaml file: %w", err)
	}

	p.conv = conv
	p.awaited = keys

	return p.parse(data)
}

func (p *Provider) parse(data []byte) (found, unknown map[string]string, err error) {
	var settings map[string]any
	if err = yaml.Unmarshal(data, &settings); err != nil {
		return nil, nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	found, unknown = p.flatten(settings)

	return found, unknown, nil
}

func (p *Provider) flatten(settings map[string]any) (found, unknown map[string]string) {
	found, unknown = make(map[string]string), make(map[string]string)

	p.flattenDFS(settings, "", found, unknown)

	return found, unknown
}

func (p *Provider) flattenDFS(m map[string]any, prefix string, found, unknown map[string]string) {
	for k, v := range m {
		newKey := k
		if prefix != "" {
			newKey = prefix + "." + k
		}

		if p.awaited[newKey] {
			found[newKey] = p.conv(v)

			continue
		}

		if subMap, ok := v.(map[string]interface{}); ok {
			p.flattenDFS(subMap, newKey, found, unknown)

			continue
		}

		unknown[newKey] = p.conv(v)
	}
}

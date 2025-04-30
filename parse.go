package zerocfg

import "fmt"

type Parser interface {
	Type() string
	Parse(awaited map[string]bool, conv func(any) string) (found, unknown map[string]string, err error)
}

func Parse(ps ...Parser) error {
	c.parsers = append(c.parsers, ps...)
	awaited := c.awaited()

	for _, p := range c.parsers {
		found, _, err := p.Parse(awaited, ToString)
		if err != nil {
			return fmt.Errorf("parse %q: %w", p.Type(), err)
		}

		err = c.applyParser(found)
		if err != nil {
			return fmt.Errorf("apply %q: %w", p.Type(), err)
		}

	}

	return nil
}

func (c *config) applyParser(vs map[string]string) error {
	for k, v := range vs {
		err := c.set(k, v)
		if err != nil {
			return fmt.Errorf("set key=%q: %w", k, err)
		}
	}

	return nil
}

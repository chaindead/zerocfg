package zfg

import "fmt"

type Parser interface {
	Type() string
	Parse() (map[string]string, error)
}

func Parse(ps ...Parser) error {
	c.parsers = append(c.parsers, ps...)

	for _, p := range c.parsers {
		vs, err := p.Parse()
		if err != nil {
			return fmt.Errorf("parse %q: %w", p.Type(), err)
		}

		err = c.applyParser(vs)
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

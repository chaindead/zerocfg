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
			return fmt.Errorf("parse %q: %v", p.Type(), err)
		}

		for k, v := range vs {
			err = c.set(k, v)
			if err != nil {
				return fmt.Errorf("set key=%q by source=%q: %v", k, p.Type(), err)
			}
		}
	}

	return nil
}

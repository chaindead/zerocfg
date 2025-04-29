package zfg

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func val[T any](v T, create func(T, *T) Value) Value {
	p := new(T)
	*p = v

	return create(v, p)
}

func Test_ConfigOk(t *testing.T) {
	const (
		name   = "name_key"
		alias  = "alias_key"
		desc   = "description"
		num    = 10
		prefix = "prefix"
	)

	tests := []struct {
		name   string
		setup  func()
		source map[string]string
		expect *config
	}{
		{
			name: "default",
			setup: func() {
				Int(name, num, desc)
				return
			},
			source: map[string]string{},
			expect: &config{
				vs: map[string]*Node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
					},
				},
			},
		},
		{
			name: "override default",
			setup: func() {
				Int(name, 0, desc)
				return
			},
			source: map[string]string{
				name: strconv.Itoa(num),
			},
			expect: &config{
				vs: map[string]*Node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						fromSource:  true,
					},
				},
			},
		},
		{
			name: "alias",
			setup: func() {
				Int(name, 0, desc, Alias(alias))
				return
			},
			source: map[string]string{
				name: strconv.Itoa(num),
			},
			expect: &config{
				vs: map[string]*Node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						Aliases:     []string{alias},
						fromSource:  true,
					},
				},
				aliases: map[string]string{
					alias: name,
				},
			},
		},
		{
			name: "override alias",
			setup: func() {
				Int(name, 0, desc, Alias(alias))
				return
			},
			source: map[string]string{
				alias: strconv.Itoa(num),
			},
			expect: &config{
				vs: map[string]*Node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						Aliases:     []string{alias},
						fromSource:  true,
					},
				},
				aliases: map[string]string{
					alias: name,
				},
			},
		},
		{
			name: "group",
			setup: func() {
				g := NewGroup(prefix)
				Int(name, 0, desc, Group(g))
				return
			},
			source: map[string]string{
				prefix + "." + name: strconv.Itoa(num),
			},
			expect: &config{
				vs: map[string]*Node{
					prefix + "." + name: {
						Name:        prefix + "." + name,
						Description: desc,
						Value:       val(num, newIntValue),
						fromSource:  true,
					},
				},
			},
		},
		{
			name: "option group",
			setup: func() {
				g := NewOptions(Secret())
				Int(name, 0, desc, Group(g))
				return
			},
			source: map[string]string{
				name: strconv.Itoa(num),
			},
			expect: &config{
				vs: map[string]*Node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						fromSource:  true,
						isSecret:    true,
					},
				},
			},
		},
		{
			name: "secret",
			setup: func() {
				Int(name, num, desc, Secret())
				return
			},
			expect: &config{
				vs: map[string]*Node{
					name: {
						Name:        name,
						Description: desc,
						Value:       val(num, newIntValue),
						isSecret:    true,
					},
				},
			},
		},
	}

	setConfig := func(expect *config) {
		c.parsers = nil
		if expect.vs == nil {
			expect.vs = make(map[string]*Node)
		}

		if expect.aliases == nil {
			expect.aliases = make(map[string]string)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c = defaultConfig()
			tt.setup()

			err := c.applyParser(tt.source)
			require.NoError(t, err)

			setConfig(tt.expect)
			require.EqualValues(t, tt.expect, c)
		})
	}
}

func Test_ConfigError(t *testing.T) {
	const (
		name    = "name_key"
		desc    = "description"
		unknown = "wrong_name"
	)

	tests := []struct {
		name    string
		setup   func()
		source  map[string]string
		err     error
		isPanic bool
	}{
		{
			name: "unknown option",
			setup: func() {
				Int(name, 0, desc)
				return
			},
			source: map[string]string{
				unknown: strconv.Itoa(0),
			},
			err: ErrNoSuchKey,
		},
		{
			name: "alias same as existing key",
			setup: func() {
				Int(name, 0, desc)
				Str(name+"1", "", desc, Alias(name))
				return
			},
			err:     fmt.Errorf("key=%q: %w", name, ErrCollidingAlias),
			isPanic: true,
		},
		{
			name: "key same as existing key",
			setup: func() {
				Int(name, 0, desc)
				Str(name, "", desc)
				return
			},
			err:     fmt.Errorf("key=%q: %w", name, ErrDuplicateKey),
			isPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := func() error {
				c = defaultConfig()
				tt.setup()

				return c.applyParser(tt.source)
			}

			if tt.isPanic {
				require.PanicsWithError(t, tt.err.Error(), func() {
					_ = fn()
				})
			} else {
				require.ErrorIs(t, fn(), tt.err)
			}
		})
	}
}

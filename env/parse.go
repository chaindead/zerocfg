package env

import (
	"os"
	"regexp"
	"strings"

	denv "github.com/joho/godotenv"
)

var cleanRe = regexp.MustCompile(`[^A-Za-z0-9.]+`)

type Opt func(*Provider)

// WithPrefix returns an Opt that sets the prefix for environment variable names in the Provider.
func WithPrefix(prefix string) Opt {
	return func(p *Provider) {
		p.prefix = prefix
	}
}

func WithPath(path *string) Opt {
	return func(p *Provider) {
		p.path = path
	}
}

// Provider parses environment variables for configuration.
type Provider struct {
	// Prefix to prepend to all environment variable names.
	prefix string
	path   *string
}

// New creates a new Provider with the provided options.
func New(opts ...Opt) *Provider {
	p := &Provider{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Type returns the type name of the parser.
func (p Provider) Type() string {
	return "env"
}

func (p Provider) key(s string) string {
	if p.prefix != "" {
		return p.prefix + "." + s
	}

	return s
}

// Parse reads environment variables matching the awaited keys and returns found values.
func (p Provider) Provide(awaited map[string]bool, _ func(any) string) (found, unknown map[string]string, err error) {
	if p.path != nil {
		err = denv.Load(*p.path)
		if err != nil {
			return nil, nil, err
		}
	}

	keys := make(map[string]string, len(awaited))
	for k := range awaited {
		keys[k] = toENV(p.key(k))
	}

	found = make(map[string]string)
	for original, formatted := range keys {
		v, ok := os.LookupEnv(formatted)
		if !ok {
			continue
		}

		found[original] = v
	}

	return found, unknown, nil
}

// toENV transforms the input string into an uppercase, underscore-separated
// environment variable name by:
// 1. Removing all characters except letters, digits, and dots.
// 2. Converting to uppercase.
// 3. Replacing dots with underscores.
func toENV(s string) string {
	cleaned := cleanRe.ReplaceAllString(s, "")
	upper := strings.ToUpper(cleaned)
	envName := strings.ReplaceAll(upper, ".", "_")

	return envName
}

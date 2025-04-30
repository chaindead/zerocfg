package env

import (
	"os"
	"regexp"
	"strings"
)

var cleanRe = regexp.MustCompile(`[^A-Za-z0-9.]+`)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p Parser) Type() string {
	return "env"
}

func (p Parser) Parse(awaited map[string]bool, _ func(any) string) (found, unknown map[string]string, err error) {
	keys := make(map[string]string, len(awaited))
	for k := range awaited {
		keys[k] = toENV(k)
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

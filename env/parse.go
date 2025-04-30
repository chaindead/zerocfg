package env

import (
	"os"
	"strings"
)

var r = strings.NewReplacer(".", "_", // separator
	// remove
	"-", "",
	"_", "",
)

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

func toENV(key string) string {
	key = strings.ToUpper(key)

	return r.Replace(key)
}

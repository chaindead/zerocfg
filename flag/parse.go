package flag

import (
	"os"
	"strings"
)

type Parser struct{}

func New() Parser {
	return Parser{}
}

func (Parser) Type() string {
	return "flag"
}

func (Parser) Parse() (map[string]string, error) {
	args := os.Args[1:]

	return parse(args)
}

func parse(args []string) (map[string]string, error) {
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		var name string
		if strings.HasPrefix(arg, "-") {
			name = arg[1:]
		}

		if strings.HasPrefix(arg, "--") {
			name = arg[2:]
		}

		if name == "" {
			continue
		}

		var value string
		if i+1 < len(args) && len(args[i+1]) > 0 && args[i+1][0] != '-' {
			value = args[i+1]
			i++
		}

		flags[name] = value
	}

	return flags, nil
}

package flag

import (
	"os"
	"strings"
)

var conv func(any) string

type Parser struct{}

func New() Parser {
	return Parser{}
}

func (Parser) Type() string {
	return "flag"
}

func (Parser) Parse(awaited map[string]bool, c func(any) string) (found, unknown map[string]string, err error) {
	args := os.Args[1:]
	conv = c

	found, unknown = parse(awaited, args)
	return
}

func parse(awaited map[string]bool, args []string) (found, unknown map[string]string) {
	found, unknown = make(map[string]string), make(map[string]string)

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

		if _, ok := awaited[name]; ok {
			found[name] = conv(value)
		} else {
			unknown[name] = conv(value)
		}
	}

	return
}

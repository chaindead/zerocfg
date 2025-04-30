# Zero Effort Configuration

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/chaindead/zerocfg)
[![Go Report Card](https://goreportcard.com/badge/github.com/chaindead/zerocfg)](https://goreportcard.com/report/github.com/chaindead/zerocfg)
[![Codecov](https://codecov.io/gh/chaindead/zerocfg/branch/main/graph/badge.svg)](https://codecov.io/gh/chaindead/zerocfg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![visitors](https://api.visitorbadge.io/api/visitors?path=github.com%2Fchaindead%2Fzerocfg&label=Visitors&countColor=%23263759&style=plastic&labelStyle=none)](https://visitorbadge.io/status?path=github.com%2Fchaindead%2Fzerocfg)

`zerocfg` is a Go package that provides fast and simple configuration management with support for multiple sources. It's designed to minimize configuration boilerplate while maintaining flexibility and ease of use.

## Features

- 🚀 Multiple configuration sources (flags, environment variables, YAML)
- 💪 Strong type safety with compile-time checks
- 🎯 Priority-based value resolution
- 🛠 Simple and intuitive API
- 📝 Automatic documentation generation
- 🔒 Built-in support for secret values

## Installation

```bash
go get -u github.com/chaindead/zerocfg
```

## Quick Start

Here's a complete example showing how to use `zerocfg`:

```go
package main

import (
	"fmt"

	zfg "github.com/chaindead/zerocfg"
	"github.com/chaindead/zerocfg/env"
	"github.com/chaindead/zerocfg/yaml"
)

var (
	// Configuration variables
	path     = zfg.Str("config.path", "", "path to yaml conf file", zfg.Alias("c"))
	ip       = zfg.IP("db.ip", "127.0.0.1", "database location")
	port     = zfg.Uint("db.port", 5678, "database port")
	username = zfg.Str("db.user", "guest", "user of database")
	password = zfg.Str("db.password", "qwerty", "password for user", zfg.Secret())
)

func main() {
	// Initialize configuration with multiple sources
	err := zfg.Parse(
		env.New(),
		yaml.New(path),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connect to %s:%d creds=%s:%s\n", *ip, *port, *username, *password)
	// OUTPUT: Connect to 127.0.0.1:5678 creds=guest:qwerty

	fmt.Println(zfg.Configuration())
	// CMD: go run ./... -c test.yaml
	// OUTPUT:
	//  config.path = test.yaml      (path to yaml conf file)
	//  db.ip       = 127.0.0.1      (database location)
	//  db.password = <secret>       (password for user)
	//  db.port     = 5678           (database port)
	//  db.user     = guest          (user of database)
}
```

## Configuration Sources

The configuration system follows a strict priority hierarchy:

1. Command-line flags (always highest priority, enabled by default)
2. Optional parsers in order of addition (first added = higher priority)
3. Default values (lowest priority)

For example, if you initialize configuration like this:
```go
zfg.Parse(
    env.New(),      // Second highest priority (after cli flags)
    yaml.New(path), // Third highest priority
)
```

The final value resolution order will be:
1. Command-line flags (if provided)
2. Environment variables (if env parser added)
3. YAML configuration (if yaml parser added)
4. Default values

Important notes:
- Lower priority sources cannot override values from higher priority sources
- All parsers except flags are optional
- Parser priority is determined by the order in `Parse()` function
- Values not found in higher priority sources fall back to lower priority sources

### Environment Variables

Environment variables are automatically transformed from the configuration key format:

| Config Key | Environment Variable | Note |
|------------|---------------------|------|
| db.user | DB_USER | Basic transformation |
| app.api.key | APP_API_KEY | Multi-level path |
| camelCase.value | CAMELCASE_VALUE | CamelCase handling |
| api-key.secret | APIKEY_SECRET | Dashes removed |
| under_score.value | UNDERSCORE_VALUE | Underscores removed |

The transformation rules:
1. Remove special characters (except letters, digits, and dots)
2. Replace dots with underscores
3. Convert to uppercase

## Documentation

For detailed documentation and advanced usage examples, visit our [Godoc page](https://godoc.org/github.com/chaindead/zerocfg).

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=chaindead/zerocfg&type=Date)](https://www.star-history.com/#chaindead/zerocfg&Date)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
# Zero Effort Configuration

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/chaindead/zerocfg) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/chaindead/zerocfg/main/LICENSE) [![codecov](https://codecov.io/gh/chaindead/zerocfg/branch/main/graph/badge.svg)](https://codecov.io/gh/chaindead/zerocfg)


The zerocfg package provides a fast and simple configuration supporting different sources.

Zerocfg's API is designed to provide both a great developer experience:
 - reducing configuration boilerplate
 - enforcing simplicity
 - provide cool features out of the box

## Installation

```bash
go get -u github.com/chaindead/zerocfg
```

## Getting Started

### Usage

```go
package main

import (
	"fmt"

	zfg "github.com/chaindead/zerocfg"
	"github.com/chaindead/zerocfg/env"
	"github.com/chaindead/zerocfg/yaml"
)

var (
	path = zfg.Str("config.path", "", "path to yaml conf file", zfg.Alias("c"))

	ip       = zfg.IP("db.ip", "127.0.0.1", "database location")
	port     = zfg.Uint("db.port", 5678, "database port")
	username = zfg.Str("db.user", "guest", "user of database")
	password = zfg.Str("db.password", "qwerty", "password for user", zfg.Secret())
)

func main() {
	// Add optional parsers with priority:
	// - flag(enabled by default)
	// - env
	// - yaml
	// Is this case most prioritised value source is flag(always), then env, then yaml,
	// So value priority is:
	// - flag passed to app
	// - parsers in order passed to Parse func (first highest)
	// - default value
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
	// OUTPUT:
	//  config.path = test/test.yaml (path to yaml conf file)
	//  db.ip       = 127.0.0.1      (database location)
	//  db.password = <secret>       (password for user)
	//  db.port     = 5678           (database port)
	//  db.user     = guest          (user of database)
}
```

### Env source

Keys for env source doest match directly, for example: 

```go
username = zfg.Str("db.user", "guest", "user of database")
```

will be parsed from env `DB_USER`, transforming key string into an uppercase, underscore-separated environment variable name by:
1. Removing all characters except letters, digits, and dots.
2. Converting to uppercase.
3. Replacing dots with underscores.

Examples:
* `db.user` -> `DB_USER`
* `camelCase.da-sh.under_wear` -> `CAMELCASE_DASH_UNDERWEAR`
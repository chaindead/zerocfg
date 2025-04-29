package main

import (
	"log"
	"zfg"
	"zfg/yaml"
)

var (
	fcmd       = zfg.Int("from.cmd", 10, "description", zfg.Alias("cmd"))
	fyaml      = zfg.Int("from.yaml", 10, "description", zfg.Alias("y"))
	configPath = zfg.String("config.path", "", "path to config file", zfg.Alias("c"))
	sl         = zfg.Strings("ss", nil, "strs")
	nl         = zfg.Ints("numbers", []int{1, 2}, "list of numbers", zfg.Alias("nl"))
)

func main() {
	err := zfg.Parse(
		yaml.New(configPath),
	)
	if err != nil {
		panic(err)
	}

	log.Println(zfg.Configuration())
}

package main

import (
	"fmt"

	confv2 "github.com/chhz0/goose/conf/v2"
)

type Config struct {
	App    string `json:"app" yaml:"app"`
	Server Server `json:"server" yaml:"server"`
}

type Server struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

func main() {
	confPtr := &Config{}
	conf, err := confv2.Init().
		WithSet("high.priority.value", "setValue").
		WithArgs(map[string]any{"conf.arg": "argValue"}).
		WithEnvPrefix("GOOSE").
		WithConfigFile("example", "yaml", ".", "./config").
		WithDefault("time", 30).
		WithUnmarshal(confPtr).
		Loading()

	if err != nil {
		fmt.Println(err)
	}
	keysmap := conf.AllSettings()
	fmt.Println(keysmap)
}

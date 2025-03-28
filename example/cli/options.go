package main

import (
	"github.com/chhz0/goose/cli"
	"github.com/spf13/pflag"
)

type RootOption struct {
	AppName string `mapstructure:"app_name" yaml:"app_name"`
	Version string `mapstructure:"version" yaml:"version"`
	Source  string `mapstructure:"source" yaml:"source"`
}

func (r *RootOption) flags() *cli.FlagSet {
	var local = func(pfs *pflag.FlagSet) {
		pfs.StringVarP(&r.Version, "version", "v", "", "app version")
		pfs.StringVarP(&r.AppName, "app_name", "a", "", "app name")
	}

	var persistent = func(pfs *pflag.FlagSet) {
		pfs.StringVarP(&r.Source, "source", "s", "", "app source")
	}

	return &cli.FlagSet{
		Local:      local,
		Persistent: persistent,
	}
}

type PrintOption struct {
	Print string `mapstructure:"print"`
	From  string `mapstructure:"from"`
}

func (p *PrintOption) flags() *cli.FlagSet {
	var local = func(pfs *pflag.FlagSet) {
		pfs.StringVarP(&p.Print, "print", "p", "print", "print")
		pfs.StringVarP(&p.From, "from", "f", "from", "from")
	}

	return &cli.FlagSet{
		Local: local,
	}
}

type EchoOption struct {
	Echo string      `mapstructure:"echo"`
	Time TimesOption `mapstructure:"time"`
}

func (e *EchoOption) flags() *cli.FlagSet {
	var local = func(pfs *pflag.FlagSet) {
		pfs.StringVarP(&e.Echo, "echo", "e", "echo", "echo")
		e.Time.flags().Local(pfs)
	}

	return &cli.FlagSet{
		Local: local,
	}
}

type TimesOption struct {
	Time int `mapstructure:"times"`
}

func (t *TimesOption) flags() *cli.FlagSet {
	var local = func(pfs *pflag.FlagSet) {
		pfs.IntVarP(&t.Time, "times", "t", 5, "times")
	}

	return &cli.FlagSet{
		Local: local,
	}
}

package cli

import "github.com/spf13/viper"

type Configer interface {
	ConfigFile() string
	ConfigPath() []string
	V() *viper.Viper
}

type config struct {
	ptr any
}

func NewConfig() *config {

	return &config{}
}

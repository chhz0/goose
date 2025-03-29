package confv1

import (
	"testing"

	"github.com/spf13/pflag"
)

// todo: 优化测试用例

type Config struct {
	App    string `json:"app" yaml:"app"`
	Server Server `json:"server" yaml:"server"`
}

type Server struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

func Test_config(t *testing.T) {
	// flags := pflag.NewFlagSet("vconfig_test", pflag.ContinueOnError)
	// flags.String("app", "vconfig_flag", "app name")
	// flags.String("server.host", "flag::127.0.0.1", "Server host")
	// flags.String("server.port", "2222", "Server port")

	// t.Setenv("VCONFIG_APP", "vconfig_env")
	// t.Setenv("VCONFIG_SERVER_Host", "env::127.0.0.1")
	// t.Setenv("VCONFIG_SERVER_PORT", "3333")

	config := NewWith(
		// 设置值
		// WithSets(map[string]any{
		// 	"app": "vconfig_set",
		// 	"server": map[string]any{
		// 		"host": "set::127.0.0.1",
		// 		"port": "9090",
		// 	},
		// }),
		// // 默认值设置
		// WithDefaults(map[string]any{
		// 	"app": "vconfig_default",
		// 	"server": map[string]any{
		// 		"host": "default::localhost",
		// 		"port": "8080",
		// 	},
		// }),

		// WithConfigName("config"),
		// WithConfigType("yaml"),
		// WithConfigPaths("./config", "."),

		WithConfig(&LocalConfig{
			ConfigName:  "config",
			ConfigType:  "yaml",
			ConfigPaths: []string{"./config"},
		}),
		WithDotEnv("dev", "."),

		WithEnvPrefix("VCONFIG"),
		// EnableDotEnv(true),
		// EnableFlag(flags),
	)
	config.Load()

	t.Log("all settings")
	t.Log(config.MarshalToString("json"))

	// var cfg *Config
	// _ = config.Unmarshal(&cfg)
	// t.Log("unmarshal", cfg)

	// config.Watcher(func() {})
}

func Test_VConfig_Set(t *testing.T) {
	config := NewWith(
		WithSets(map[string]any{
			"app": "vconfig_set",
			"server": map[string]any{
				"host": "set::127.0.0.1",
				"port": "9090",
			},
		}),
	)

	config.Set("set", "vconfig_set_new")

	config.Load()

	t.Log("all settings")
	t.Log(config.AllSettings())
}

func Test_VConfig_Flag(t *testing.T) {
	flags := pflag.NewFlagSet("vconfig_test", pflag.ContinueOnError)
	flags.String("app", "vconfig_flag", "app name")
	flags.String("server.host", "flag::127.0.0.1", "Server host")
	flags.String("server.port", "2222", "Server port")

	flag2 := pflag.NewFlagSet("vconfig_test2", pflag.ContinueOnError)
	flag2.String("app1", "vconfig_flag2", "app name")
	flag2.String("server.host1", "flag2::127.0.0.1", "Server host")
	flag2.String("server.port1", "2222", "Server port")

	config := NewWith(EnableFlag(flags, flag2))

	config.BindPFlag(map[string]*pflag.Flag{
		"bind": flags.Lookup("app"),
	})
	config.Load()
	t.Log("all settings")
	t.Log(config.AllSettings())
}

func Test_VConfig_Env(t *testing.T) {
	t.Setenv("VCONFIG_APP", "vconfig_env")
	t.Setenv("VCONFIG_SERVER_HOST", "env::127.0.0.1")
	t.Setenv("VCONFIG_SERVER_PORT", "3333")

	config := NewWith(
		WithEnv(&Env{
			Binds:  []string{"app", "server.host", "server.port"},
			Prefix: "VCONFIG",
		}),
	)
	config.Load()
	// config.BindEnvs("app")
	// config.BindEnvs("server.host")
	// config.BindEnvs("server.port")
	t.Log("all settings")
	t.Log(config.MarshalToString("json"))
}

func Test_VConfig_DotEnv(t *testing.T) {
	config := NewWith(
		WithDotEnv("dev", "."),
	)
	config.Load()
	t.Log("all settings")
	t.Log(config.MarshalToString("json"))
}

func Test_VConfig_Config(t *testing.T) {
	config := NewWith(
		WithConfig(&LocalConfig{
			ConfigName:  "config",
			ConfigType:  "yaml",
			ConfigPaths: []string{"./config"},
		}),
	)
	config.Load()
	t.Log("all settings")
	t.Log(config.MarshalToString("json"))
}

func Test_VConfig_Remote(t *testing.T) {
	// TODO: to do
}

func Test_VConfig_Default(t *testing.T) {
	config := NewWith(
		WithDefaults(map[string]any{
			"app": "vconfig_default",
			"server": map[string]any{
				"host": "default::localhost",
				"port": "8080",
			},
		}),
	)
	config.Load()
	t.Log("all settings")
	t.Log(config.MarshalToString("json"))
}

func Test_VConfig_KeyValue(t *testing.T) {
	// TODO: to do
}

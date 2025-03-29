package confv1

import (
	"io"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type RemoteProvider struct {
	Provider string
	Endpoint string
	Path     string
	Type     string
}

type Env struct {
	Binds       []string // 环境变量命
	Prefix      string   // 环境变量前缀
	KeyReplacer *strings.Replacer
	// TODO: allow empty env
}

// TODO: 多配置文件来源
type LocalConfig struct {
	ConfigName  string    // 配置文件名
	ConfigType  string    // 配置文件类型
	ConfigPaths []string  // 配置文件路径
	ConfigIO    io.Reader // 配置读取 IO
}

type Options struct {
	Sets     map[string]any
	Defaults map[string]any

	Config *LocalConfig
	DotEnv *LocalConfig

	Env *Env

	Flags []*pflag.FlagSet // flags

	// UnmarshalPtr 反序列化对象, 必须是 指针
	// 如果提供了 UnmarshalPtr 且开启了Watcher，在配置文件更新时自动反序列化
	UnmarshalPtr any

	RemoteS             struct{}
	Remote              *RemoteProvider
	RemoteWatch         bool
	RemoteWatchInterval time.Duration

	EnableEnv    bool // 是否开启环境变量
	EnableFlag   bool // 是否使用flag
	EnableRemote bool // 是否开启远程配置中心
}

func WithSets(sets map[string]any) func(*Options) {
	return func(o *Options) {
		o.Sets = sets
	}
}
func WithDefaults(defaluts map[string]any) func(*Options) {
	return func(o *Options) {
		o.Defaults = defaluts
	}
}

func WithDotEnv(mode string, path ...string) func(*Options) {
	return func(o *Options) {
		o.DotEnv = &LocalConfig{
			ConfigName:  mode,
			ConfigType:  "env",
			ConfigPaths: path,
		}
	}
}

func WithConfig(local *LocalConfig) func(*Options) {
	return func(o *Options) {
		o.Config = local
	}
}

func WithConfigName(name string) func(*Options) {
	return func(o *Options) {
		o.Config.ConfigName = name
	}
}

func WithConfigType(configType string) func(*Options) {
	return func(o *Options) {
		o.Config.ConfigType = configType
	}
}

func WithConfigPaths(paths ...string) func(*Options) {
	return func(o *Options) {
		o.Config.ConfigPaths = append(o.Config.ConfigPaths, paths...)
	}
}

func WithUnmarshal(ptr any) func(*Options) {
	return func(o *Options) {
		o.UnmarshalPtr = ptr
	}
}

// WithEnv 允许设置环境变量, 如果使用 WithEnv ， 必须传入的 Env.KeyReplacer
func WithEnv(env *Env) func(*Options) {
	return func(o *Options) {
		if env.KeyReplacer == nil {
			env.KeyReplacer = defaultKeyReplacer()
		}
		o.Env = env
	}
}

func WithEnvBinds(binds ...string) func(*Options) {
	return func(o *Options) {
		o.Env.Binds = append(o.Env.Binds, binds...)
	}
}

func WithEnvPrefix(prefix string) func(*Options) {
	return func(o *Options) {
		o.Env.Prefix = prefix
	}
}

func WithEnvKeyReplacer(replacer *strings.Replacer) func(*Options) {
	return func(o *Options) {
		o.Env.KeyReplacer = replacer
	}
}

func WithRemote(remote *RemoteProvider) func(*Options) {
	return func(o *Options) {
		o.Remote = remote
	}
}

func EnableEnv(enable bool) func(*Options) {
	return func(o *Options) {
		o.EnableEnv = enable
	}
}

func EnableFlag(flags ...*pflag.FlagSet) func(*Options) {
	return func(o *Options) {
		o.EnableFlag = true
		o.Flags = append(o.Flags, flags...)
	}
}

func EnableRemote(enable bool) func(*Options) {
	return func(o *Options) {
		o.EnableRemote = enable
	}
}

func EnableRemoteWatch(enable bool) func(*Options) {
	return func(o *Options) {
		o.RemoteWatch = enable
	}
}

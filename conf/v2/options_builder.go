package confv2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ErrConfigNotFound = errors.New("config file not found")
	ErrConfigRead     = errors.New("config file read error")
	ErrDotEnvNotFound = errors.New("dotenv file not found")
	ErrDotEnvRead     = errors.New("dotenv file read error")
	ErrReaderIO       = errors.New("reader io error")
	ErrRemoteConfig   = errors.New("remote config error")
	ErrUnmarshal      = errors.New("unmarshal error")
)

type Options struct {
	sets map[string]any
	args map[string]any

	envPrefix   string
	envReplacer *strings.Replacer
	envBinds    []string

	dotEnv     *FileConfig
	configFile *FileConfig
	remote     *RemoteConfig

	defaults map[string]any

	flags []*pflag.FlagSet

	unmarshalTo   any
	watching      bool
	watchRemote   bool
	watchInterval time.Duration

	errReadHandler func(err error) error
}

type FileConfig struct {
	name  string
	typ   string
	paths []string
	io    io.Reader
}

type RemoteConfig struct {
	provider string
	endpoint string
	path     string
	typ      string
}

type ConfigBuilder struct {
	opts Options
	err  error
}

func Init() *ConfigBuilder {
	return &ConfigBuilder{
		opts: Options{
			sets:           make(map[string]any),
			args:           make(map[string]any),
			defaults:       make(map[string]any),
			envReplacer:    strings.NewReplacer(".", "_", "-", "_"),
			errReadHandler: func(err error) error { return err },
		},
	}
}

func (b *ConfigBuilder) Loading() (*Config, error) {
	if b.err != nil {
		return nil, b.err
	}

	v := viper.New()
	c := &Config{
		v:    v,
		opts: b.opts,
	}

	c.setDefault()
	if err := c.loadConfigFile(); err != nil {
		return nil, err
	}
	if err := c.loadDotEnv(); err != nil {
		return nil, err
	}
	c.setupEnv()
	c.bindPFlags()
	c.setArgs()
	c.set()

	if c.opts.watching {
		c.watchConfig()
	}
	if c.opts.watchRemote && c.opts.remote != nil {
		go c.watchRemote(context.Background())
	}

	return c, nil
}

func (b *ConfigBuilder) MustLoading() *Config {
	cfg, err := b.Loading()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	return cfg
}

func (b *ConfigBuilder) WithSet(key string, value any) *ConfigBuilder {
	b.opts.sets[key] = value
	return b
}

func (b *ConfigBuilder) WithArgs(args map[string]any) *ConfigBuilder {
	for k, v := range args {
		b.opts.args[k] = v
	}

	return b
}

func (b *ConfigBuilder) WithEnvPrefix(prefix string) *ConfigBuilder {
	b.opts.envPrefix = prefix
	return b
}

func (b *ConfigBuilder) WithEnvKeyReplacer(replacer *strings.Replacer) *ConfigBuilder {
	b.opts.envReplacer = replacer
	return b
}

func (b *ConfigBuilder) WithEnvBind(keys ...string) *ConfigBuilder {
	b.opts.envBinds = append(b.opts.envBinds, keys...)
	return b
}

func (b *ConfigBuilder) WithDotEnv(name string, paths ...string) *ConfigBuilder {
	b.opts.dotEnv = &FileConfig{
		name:  name,
		typ:   "env",
		paths: paths,
	}
	return b
}

func (b *ConfigBuilder) WithConfigFile(name, typ string, paths ...string) *ConfigBuilder {
	b.opts.configFile = &FileConfig{
		name:  name,
		typ:   typ,
		paths: paths,
	}
	return b
}

func (b *ConfigBuilder) WithConfigReader(r io.Reader, typ string) *ConfigBuilder {
	if b.opts.configFile == nil {
		b.opts.configFile = &FileConfig{}
	}
	b.opts.configFile.io = r
	b.opts.configFile.typ = typ
	return b
}

func (b *ConfigBuilder) WithDefault(key string, value any) *ConfigBuilder {
	b.opts.defaults[key] = value
	return b
}

func (b *ConfigBuilder) WithRemote(provider, endpoint, path, typ string) *ConfigBuilder {
	b.opts.remote = &RemoteConfig{
		provider: provider,
		endpoint: endpoint,
		path:     path,
		typ:      typ,
	}
	return b
}

func (b *ConfigBuilder) WithFlags(flags ...*pflag.FlagSet) *ConfigBuilder {
	b.opts.flags = append(b.opts.flags, flags...)
	return b
}

func (b *ConfigBuilder) WithUnmarshal(target any) *ConfigBuilder {
	b.opts.unmarshalTo = target
	return b
}

func (b *ConfigBuilder) WithWatch(enable bool) *ConfigBuilder {
	b.opts.watching = enable
	return b
}

func (b *ConfigBuilder) WithRemoteWatch(enable bool, interval time.Duration) *ConfigBuilder {
	b.opts.watchRemote = enable
	b.opts.watchInterval = interval
	return b
}

func (b *ConfigBuilder) WithErrReadHandler(handler func(err error) error) *ConfigBuilder {
	b.opts.errReadHandler = handler
	return b
}

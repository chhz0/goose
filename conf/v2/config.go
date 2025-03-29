package confv2

import (
	"context"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	v    *viper.Viper
	opts Options
	mu   sync.RWMutex
}

func (c *Config) setDefault() {
	for k, v := range c.opts.defaults {
		c.v.SetDefault(k, v)
	}
}

func (c *Config) loadConfigFile() error {
	if c.opts.configFile == nil {
		return nil
	}

	if c.opts.configFile.io != nil {
		c.v.SetConfigType(c.opts.configFile.typ)
		if err := c.v.ReadConfig(c.opts.configFile.io); err != nil {
			return c.opts.errReadHandler(ErrReaderIO)
		}
		return nil
	}

	c.v.SetConfigName(c.opts.configFile.name)
	c.v.SetConfigType(c.opts.configFile.typ)
	for _, path := range c.opts.configFile.paths {
		c.v.AddConfigPath(path)
	}

	if err := c.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrConfigNotFound
		}
		return c.opts.errReadHandler(ErrConfigRead)
	}

	return nil
}

func (c *Config) loadDotEnv() error {
	if c.opts.dotEnv == nil {
		return nil
	}

	v := viper.New()
	v.SetConfigName(c.opts.dotEnv.name)
	v.SetConfigType(c.opts.dotEnv.typ)
	for _, path := range c.opts.dotEnv.paths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrDotEnvNotFound
		}
		return c.opts.errReadHandler(ErrDotEnvRead)
	}

	return c.v.MergeConfigMap(v.AllSettings())
}

func (c *Config) setupEnv() {
	if c.opts.envPrefix != "" {
		c.v.SetEnvPrefix(c.opts.envPrefix)
	}
	if c.opts.envReplacer != nil {
		c.v.SetEnvKeyReplacer(c.opts.envReplacer)
	}
	c.v.AutomaticEnv()
	_ = c.v.BindEnv(c.opts.envBinds...)
}

func (c *Config) bindPFlags() {
	for _, fs := range c.opts.flags {
		_ = c.v.BindPFlags(fs)
	}
}

func (c *Config) setArgs() {
	for k, v := range c.opts.args {
		c.v.Set(k, v)
	}
}

func (c *Config) set() {
	for k, v := range c.opts.sets {
		c.v.Set(k, v)
	}
}

func (c *Config) watchConfig() {
	c.v.OnConfigChange(func(in fsnotify.Event) {
		c.mu.Lock()
		defer c.mu.Unlock()

		if err := c.v.ReadInConfig(); err != nil {
			_ = c.opts.errReadHandler(err)
			return
		}

		if c.opts.unmarshalTo != nil {
			if err := c.v.Unmarshal(c.opts.unmarshalTo); err != nil {
				return
			}
		}
	})
	c.v.WatchConfig()
}

func (c *Config) watchRemote(ctx context.Context) {
	if c.opts.remote == nil {
		return
	}

	ticker := time.NewTicker(c.opts.watchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			if err := c.v.WatchRemoteConfig(); err != nil {
				return
			}
			c.mu.Unlock()
		}
	}
}

func (c *Config) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.Get(key)
}

func (c *Config) Unmarshal(target any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v.Unmarshal(target)
}

func (c *Config) AllSettings() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.AllSettings()
}

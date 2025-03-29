// viper.set > args > env > config file > key/value store > default
package confv1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	ErrConfigNotFound = errors.New("config file not found")
	ErrDotEnvNotFound = errors.New("dotenv file not found")
	ErrReaderIO       = errors.New("reader io error")
	ErrInvalidType    = errors.New("invalid config type")
	ErrRemoteConfig   = errors.New("remote config error")
	ErrUnmarshal      = errors.New("unmarshal error")
	ErrUnmarshalNil   = errors.New("unmarshal nil")
)

type VConfig struct {
	v *viper.Viper

	// todo: 支持同时读取多个配置
	vps  map[string]*viper.Viper
	opts *Options
	mu   sync.RWMutex
}

// New 使用 options 模式创建配置实例
func NewWith(optFuncs ...func(*Options)) *VConfig {
	defaultOpts := &Options{
		Config: &LocalConfig{
			ConfigName:  "",
			ConfigPaths: []string{"."},
		},
		Env: &Env{
			KeyReplacer: defaultKeyReplacer(),
		},
		EnableEnv:           true,
		RemoteWatchInterval: 30 * time.Second,
	}
	for _, fn := range optFuncs {
		fn(defaultOpts)
	}

	vc := &VConfig{
		v:    viper.New(),
		vps:  make(map[string]*viper.Viper, 0),
		opts: defaultOpts,
	}

	vc.initialize()

	return vc
}

// NewInOptions 使用Options创建配置实例
// 预期：opts 必须全部配置
func New(opts *Options) *VConfig {
	vc := &VConfig{
		v:    viper.New(),
		opts: opts,
	}

	vc.initialize()

	return vc
}

func (vc *VConfig) initialize() {
	vc.setDefault()

	// 绑定 flag 参数
	if vc.opts.EnableFlag {
		vc.bindFlags()
	}

	// 加载环境变量
	if vc.opts.EnableEnv {
		vc.setupEnv()
	}

	// 加载 key/value 参数
	for key, val := range vc.opts.Sets {
		vc.v.Set(key, val)
	}
}

func (vc *VConfig) setDefault() {
	for k, v := range vc.opts.Defaults {
		vc.v.SetDefault(k, v)
	}
}

func (vc *VConfig) bindFlags() {
	for _, fs := range vc.opts.Flags {
		fs.VisitAll(func(f *pflag.Flag) {
			if err := vc.v.BindPFlag(f.Name, f); err != nil {
				log.Printf("failed to bind flag %s: %v", f.Name, err)
			}
		})
	}
}

func (vc *VConfig) setupEnv() {
	vc.v.AutomaticEnv()
	if vc.opts.Env.Prefix != "" {
		vc.v.SetEnvPrefix(vc.opts.Env.Prefix)
	}
	if vc.opts.Env.Binds != nil {
		for _, env := range vc.opts.Env.Binds {
			_ = vc.v.BindEnv(env)
		}
	}
	if vc.opts.Env.KeyReplacer != nil {
		vc.v.SetEnvKeyReplacer(vc.opts.Env.KeyReplacer)
	}

	vc.setInRead("config")

	// todo: 支持 dotenv的读取
	// if vc.opts.DotEnv != nil {
	// 	vc.setInRead("dotenv")
	// }
}

func (vc *VConfig) setInRead(in string) {
	use := vc.opts.Config
	if in == "dotenv" {
		use = vc.opts.DotEnv
	}

	vc.v.SetConfigName(use.ConfigName)
	vc.v.SetConfigType(use.ConfigType)
	for _, cp := range use.ConfigPaths {
		vc.v.AddConfigPath(cp)
	}
}

func (vc *VConfig) Load() {

	// 加载本地配置文件
	if err := vc.loadConfig(); err != nil && !errors.Is(err, ErrConfigNotFound) {
		log.Printf("Warning: Error loading local file: %v", err)
	}

	if vc.opts.DotEnv != nil {
		if err := vc.loadDotEnv(); err != nil && !errors.Is(err, ErrConfigNotFound) {
			log.Printf("Warning: Error loading local file: %v", err)
		}
	}

	// todo: 加载远程配置文件
	// if vc.opts.EnableRemote {
	// 	if err := vc.loadRemote(); err != nil {
	// 		log.Printf("Warning: Error loading remote config: %v", err)
	// 	}
	// }

}

func (vc *VConfig) loadConfig() error {
	if err := vc.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && vc.opts.Config.ConfigIO != nil {
			return vc.loadReaderIO()
		}
		return fmt.Errorf("config file read error: %v\n", err)
	}

	return nil
}

func (vc *VConfig) loadDotEnv() error {
	if err := vc.v.MergeInConfig(); err != nil {
		if os.IsNotExist(err) {
			return ErrDotEnvNotFound
		}
		return fmt.Errorf("dotenv file merge error: %v\n", err)
	}
	return nil
}

// func (vc *VConfig) mergeFromViper(vp *viper.Viper) error {
// 	return vc.v.MergeConfigMap(vp.AllSettings())
// }

func (vc *VConfig) loadReaderIO() error {
	if err := vc.v.ReadConfig(vc.opts.Config.ConfigIO); err != nil {
		return ErrReaderIO
	}

	return nil
}

func (vc *VConfig) loadRemote() error {
	if vc.opts.Remote == nil {
		return ErrRemoteConfig
	}

	remote := vc.opts.Remote
	if err := vc.v.AddRemoteProvider(remote.Provider, remote.Endpoint, remote.Path); err != nil {
		log.Printf("failed to remote provider: %v\n", err)
		return ErrRemoteConfig
	}

	vc.v.SetConfigType(remote.Type)
	if err := vc.v.ReadRemoteConfig(); err != nil {
		return ErrRemoteConfig
	}

	return nil
}

// Watcher 监听配置文件变化, changedFunc 将在配置文件更新并重新加载完成后调用
func (vc *VConfig) Watcher(changedFunc func()) {
	vc.enableWatch(changedFunc)
}

func (vc *VConfig) enableWatch(fn func()) {
	vc.v.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("config file changed: %v\n", in.Name)
		if err := vc.v.ReadInConfig(); err != nil {
			log.Printf("reload config file error: %v\n", err)
		}
		_ = vc.unmarshal()
		fn()
	})
	vc.v.WatchConfig()

	if vc.opts.RemoteWatch {
		go vc.watchRemote(context.Background())
	}
}

func (vc *VConfig) watchRemote(ctx context.Context) {
	ticker := time.NewTicker(vc.opts.RemoteWatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := vc.v.WatchRemoteConfig(); err != nil {
				log.Printf("reload remote config error: %v\n", err)
			}
		}
	}
}

func (vc *VConfig) Unmarshal(ptr any) error {
	if err := vc.v.Unmarshal(ptr); err != nil {
		return ErrUnmarshal
	}

	return nil
}

func (vc *VConfig) unmarshal() error {
	if vc.opts.UnmarshalPtr == nil {
		return ErrUnmarshalNil
	}
	if err := vc.v.Unmarshal(vc.opts.UnmarshalPtr); err != nil {
		return ErrUnmarshal
	}

	return nil
}

// Marshal 将vc.v.AllSettings()序列化为字符串
// 目前支持：json, yaml, toml
func (vc *VConfig) MarshalToString(marshalType string) (string, error) {
	m := vc.v.AllSettings()
	var buf []byte
	var err error
	switch marshalType {
	case "json":
		buf, err = json.Marshal(m)
	case "yaml":
		buf, err = yaml.Marshal(m)
	case "toml":
		buf, err = toml.Marshal(m)
	}
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (vc *VConfig) BindPFlag(mFlag map[string]*pflag.Flag) {
	for key, flag := range mFlag {
		_ = vc.v.BindPFlag(key, flag)
	}
}

func (vc *VConfig) BindPFlags(pfs ...*pflag.FlagSet) {
	for _, pf := range pfs {
		_ = vc.v.BindPFlags(pf)
	}
}

// BindEnvs 绑定环境变量，不同于viper.BindEnv限制一个传入的参数
// 如果想使用viper.BindEnv，请调用函数 V() 获取 *viper.Viper实例
func (vc *VConfig) BindEnvs(input string) {
	_ = vc.v.BindEnv(input)
}

func (vc *VConfig) GetEnv(key string) string {
	return vc.v.GetString(key)
}

func (vc *VConfig) Set(key string, value any) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.v.Set(key, value)
}

// Get 允许访问给定key 的value
// 如果有嵌套的key，则使用点号分隔符访问："section.key"
func (vc *VConfig) Get(key string) (any, bool) {
	if !vc.v.IsSet(key) {
		return nil, false
	}

	v := vc.v.Get(key)
	return v, true
}

func (vc *VConfig) AllSettings() map[string]any {
	return vc.v.AllSettings()
}

// V returns the viper instance
func (vc *VConfig) V() *viper.Viper {
	return vc.v
}

func defaultKeyReplacer() *strings.Replacer {
	return strings.NewReplacer(".", "_", "-", "_")
}

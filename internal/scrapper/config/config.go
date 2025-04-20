package config

import (
	"fmt"
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Env       string          `mapstructure:"env"       validate:"required,oneof=local prod"`
	Server    ServerConfig    `mapstructure:"server"    validate:"required"`
	Scheduler SchedulerConfig `mapstructure:"scheduler" validate:"required"`
	Bot       BotConfig       `mapstructure:"bot"       validate:"required"`
	Fetch     FetchBlock      `mapstructure:"fetch"     validate:"required"`
}

type FetchBlock struct {
	GitHub        FetchConfig `mapstructure:"github"        validate:"required"`
	StackOverflow FetchConfig `mapstructure:"stack_overflow" validate:"required"`
}

type FetchConfig struct {
	BaseURL   string        `mapstructure:"base_url"  validate:"required,url"`
	AuthToken string        `mapstructure:"auth_token"`
	Timeout   time.Duration `mapstructure:"timeout"   validate:"required"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"             validate:"required,hostname|ip"`
	Port            string        `mapstructure:"port"             validate:"required,numeric"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"     validate:"required"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"    validate:"required"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required"`
}

func (s *ServerConfig) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}

type SchedulerConfig struct {
	Interval time.Duration `mapstructure:"interval" validate:"required"`
	Timeout  time.Duration `mapstructure:"timeout"  validate:"required"`
}

type BotConfig struct {
	Host    string        `mapstructure:"host"    validate:"required"`
	Port    string        `mapstructure:"port"    validate:"required"`
	Timeout time.Duration `mapstructure:"timeout" validate:"required"`
}

func (b *BotConfig) BaseURL() string {
	return fmt.Sprintf("http://%s:%s", b.Host, b.Port)
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetDefault("env", "local")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if err := validator.New().Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	return &cfg, nil
}

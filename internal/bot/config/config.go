package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string         `mapstructure:"env" validate:"required,oneof=local prod"`
	Telegram TelegramConfig `mapstructure:"telegram" validate:"required"`
	Server   ServerConfig   `mapstructure:"server" validate:"required"`
	Scrapper ScrapperConfig `mapstructure:"scrapper" validate:"required"`
}

type TelegramConfig struct {
	Token string `mapstructure:"token" validate:"required"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host" validate:"required,hostname|ip"`
	Port            string        `mapstructure:"port" validate:"required,numeric"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" validate:"required"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required"`
}

func (c *ServerConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type ScrapperConfig struct {
	Host    string        `mapstructure:"host" validate:"required,hostname|ip"`
	Port    string        `mapstructure:"port" validate:"required,numeric"`
	Timeout time.Duration `mapstructure:"timeout" validate:"required"`
}

func (c *ScrapperConfig) Address() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(c.Host, c.Port))
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.BindEnv("telegram.token", "TELEGRAM_BOT_TOKEN"); err != nil {
		return nil, fmt.Errorf("bind env TELEGRAM_TOKEN: %w", err)
	}

	viper.SetDefault("env", "local")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
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

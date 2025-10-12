package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	App struct {
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"app"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	configName := detectEnvironment()

	v.SetConfigName(configName)
	v.AddConfigPath(configPath)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	log.Printf("Config loaded from %s (detected environment: %s)", v.ConfigFileUsed(), configName)
	return &cfg, nil
}

func detectEnvironment() string {

	hostname, _ := os.Hostname()

	if strings.Contains(hostname, "subscriptions") {
		return "config.docker"
	}

	if canResolve("postgres") {
		return "config.docker"
	}

	return "config.local"
}

func canResolve(host string) bool {
	_, err := net.LookupHost(host)
	return err == nil
}

func (p *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.Database.User,
		p.Database.Password,
		p.Database.Host,
		p.Database.Port,
		p.Database.Name,
		p.Database.SSLMode,
	)
}

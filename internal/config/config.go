// Package config loads application configuration from YAML files.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	JWT    JWTConfig    `yaml:"jwt"`
	Log    LogConfig    `yaml:"log"`
}

type ServerConfig struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
	Addr string `yaml:"addr"`
}

type MySQLConfig struct {
	DSN             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int64  `yaml:"expire"`
}

type LogConfig struct {
	Path           string `yaml:"path"`
	Level          string `yaml:"level"`
	KeepHours      int    `yaml:"keep_hours"`
	Filename       string `yaml:"filename"`
	AccessFilename string `yaml:"access_filename"`
	ErrorFilename  string `yaml:"error_filename"`
}

func MustLoad(path string) *Config {
	contents, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read config %q: %w", path, err))
	}

	var cfg Config
	if err := yaml.Unmarshal(contents, &cfg); err != nil {
		panic(fmt.Errorf("decode config %q: %w", path, err))
	}
	return &cfg
}

package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger      LoggerConfig     `yaml:"logger"`
	GrpcServer  GrpcServerConfig `yaml:"grpcServer"`
	HTTPServer  HTTPServerConfig `yaml:"httpServer"`
	Database    DatabaseConfig   `yaml:"database"`
	StorageType string           `yaml:"storage"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type GrpcServerConfig struct {
	Port string `yaml:"port"`
}

type HTTPServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	ReadTimeout time.Duration `yaml:"readTimeout"`
}

type DatabaseConfig struct {
	Driver    string `yaml:"driver"`
	DsnPrefix string `yaml:"dsnPrefix"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
}

func NewConfig(path string) (c *Config, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	return &config, err
}

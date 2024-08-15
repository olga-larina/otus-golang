package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger      LoggerConfig     `mapstructure:"logger"`
	GrpcServer  GrpcServerConfig `mapstructure:"grpcServer"`
	HTTPServer  HTTPServerConfig `mapstructure:"httpServer"`
	Database    DatabaseConfig   `mapstructure:"database"`
	StorageType string           `mapstructure:"storage"`
	Timezone    string           `mapstructure:"timezone"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type GrpcServerConfig struct {
	Port string `mapstructure:"port"`
}

type HTTPServerConfig struct {
	Host        string        `mapstructure:"host"`
	Port        string        `mapstructure:"port"`
	ReadTimeout time.Duration `mapstructure:"readTimeout"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	URI    string `mapstructure:"uri"`
}

func NewConfig(path string) (*Config, error) {
	parser := viper.New()
	parser.SetConfigFile(path)

	err := parser.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	for _, key := range parser.AllKeys() {
		value := parser.GetString(key)
		parser.Set(key, os.ExpandEnv(value))
	}

	var config Config
	err = parser.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	return &config, err
}

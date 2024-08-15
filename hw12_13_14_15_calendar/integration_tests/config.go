//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	Calendar CalendarConfig `mapstructure:"calendar"`
	Timezone string         `mapstructure:"timezone"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	URI    string `mapstructure:"uri"`
}

type CalendarConfig struct {
	GrpcURL          string        `mapstructure:"grpc_url"`
	HTTPUrl          string        `mapstructure:"http_url"`
	NotifyPeriod     time.Duration `mapstructure:"notifyPeriod"`
	NotifyScanPeriod time.Duration `mapstructure:"notifyScanPeriod"`
	ClearPeriod      time.Duration `mapstructure:"clearPeriod"`
	NotifyCronPeriod time.Duration `mapstructure:"notifyCronPeriod"`
	ClearCronPeriod  time.Duration `mapstructure:"clearCronPeriod"`
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

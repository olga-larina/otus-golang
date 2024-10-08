package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	Queue    QueueConfig    `mapstructure:"queue"`
	Schedule ScheduleConfig `mapstructure:"schedule"`
	Timezone string         `mapstructure:"timezone"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	URI    string `mapstructure:"uri"`
}

type QueueConfig struct {
	URI          string `mapstructure:"uri"`
	ExchangeName string `mapstructure:"exchangeName"`
	ExchangeType string `mapstructure:"exchangeType"`
	QueueName    string `mapstructure:"queueName"`
	RoutingKey   string `mapstructure:"routingKey"`
}

type ScheduleConfig struct {
	NotifyCron       string        `mapstructure:"notifyCron"`
	ClearCron        string        `mapstructure:"clearCron"`
	NotifyPeriod     time.Duration `mapstructure:"notifyPeriod"`
	NotifyScanPeriod time.Duration `mapstructure:"notifyScanPeriod"`
	ClearPeriod      time.Duration `mapstructure:"clearPeriod"`
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

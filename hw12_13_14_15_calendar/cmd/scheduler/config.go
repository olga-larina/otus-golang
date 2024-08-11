package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Queue    QueueConfig    `yaml:"queue"`
	Schedule ScheduleConfig `yaml:"schedule"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
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

type QueueConfig struct {
	URI          string `yaml:"uri"`
	ExchangeName string `yaml:"exchangeName"`
	ExchangeType string `yaml:"exchangeType"`
	QueueName    string `yaml:"queueName"`
	RoutingKey   string `yaml:"routingKey"`
}

type ScheduleConfig struct {
	NotifyCron       string        `yaml:"notifyCron"`
	ClearCron        string        `yaml:"clearCron"`
	NotifyPeriod     time.Duration `yaml:"notifyPeriod"`
	NotifyScanPeriod time.Duration `yaml:"notifyScanPeriod"`
	ClearPeriod      time.Duration `yaml:"clearPeriod"`
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

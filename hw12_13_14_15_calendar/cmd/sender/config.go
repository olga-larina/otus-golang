package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger LoggerConfig `yaml:"logger"`
	Queue  QueueConfig  `yaml:"queue"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type QueueConfig struct {
	URI          string `yaml:"uri"`
	ExchangeName string `yaml:"exchangeName"`
	ExchangeType string `yaml:"exchangeType"`
	QueueName    string `yaml:"queueName"`
	RoutingKey   string `yaml:"routingKey"`
	ConsumerTag  string `yaml:"consumerTag"`
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

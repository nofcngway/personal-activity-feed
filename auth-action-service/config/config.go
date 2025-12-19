package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	GRPC     GRPCConfig     `yaml:"grpc"`
	HTTP     HTTPConfig     `yaml:"http"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Kafka    KafkaConfig    `yaml:"kafka"`
	Security SecurityConfig `yaml:"security"`
}

type GRPCConfig struct {
	Addr string `yaml:"addr"`
}

type HTTPConfig struct {
	Addr       string `yaml:"addr"`
	SwaggerPath string `yaml:"swagger_path"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type KafkaConfig struct {
	Brokers   []string `yaml:"brokers"`
	TopicName string   `yaml:"topic_name"`
}

type SecurityConfig struct {
	SessionTTLSeconds int `yaml:"session_ttl_seconds"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &cfg, nil
}



package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Cockroach struct {
		DSN string `yaml:"dsn"`
	} `yaml:"cockroach"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group_id"`
	} `yaml:"kafka"`
	Cassandra struct {
		Hosts    []string `yaml:"hosts"`
		Keyspace string   `yaml:"keyspace"`
	} `yaml:"cassandra"`
	ClickHouse struct {
		Addr     string `yaml:"addr"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"clickhouse"`
	Server struct {
		Addr string `yaml:"addr"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func Load(path string) (Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading file: %w", err)
	}

	c, err := Parse(f)
	if err != nil {
		return Config{}, fmt.Errorf("parsing configs: %w", err)
	}

	return c, nil
}

// Parse reads the yaml data into a Config struct. It does not perform any validations on the configurations themselves.
func Parse(data []byte) (Config, error) {
	c := Config{}
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return Config{}, fmt.Errorf("parsing yaml file: %w", err)
	}
	return c, nil
}

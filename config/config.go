package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mongo      Mongo `yaml:"mongo"`
	HTTPServer `yaml:"http_server"`
	App        App `yaml:"app"`
}

type Mongo struct {
	ApplyURI string        `yaml:"apply_uri"`
	MaxConns uint64        `yaml:"max_conns"`
	MinConns uint64        `yaml:"min_conns"`
	IdleTime time.Duration `yaml:"idle_time"`
	DBName   string        `yaml:"db_name"`
}

type HTTPServer struct {
	Addr string `yaml:"addr"`
}

type App struct {
	GracefulTimeout time.Duration `yaml:"graceful_timeout"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

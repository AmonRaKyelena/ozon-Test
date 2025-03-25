package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

type RepositoryMode string

var (
	PostgresqlMode RepositoryMode = "postgresql"
	InMemoryMode   RepositoryMode = "inmemory"
)

type Config struct {
	Port           string `json:"port"`
	LogLevel       string `json:"log_level"`
	RepositoryMode string `json:"repository_mode"`
	PsqlInfo       string `json:"psql_info"`
}

func NewConfig() (*Config, error) {
	configPath := flag.String("config", "config.json", "path to the configuration file")
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %q: %w", *configPath, err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode JSON config: %w", err)
	}

	if cfg.Port == "" {
		return nil, errors.New("port must not be empty")
	}
	if cfg.LogLevel == "" {
		return nil, errors.New("log_level must not be empty")
	}
	if cfg.RepositoryMode == "" {
		return nil, errors.New("repository_mode must not be empty")
	}
	if cfg.RepositoryMode == string(PostgresqlMode) && cfg.PsqlInfo == "" {
		return nil, errors.New("psql_info must not be empty in this repository mode")
	}

	return &cfg, nil
}

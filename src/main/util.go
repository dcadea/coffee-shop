package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func loadConfig(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return config, nil
}

func getDBHost() string {
	dbHost := os.Getenv("DB_HOST")

	if len(dbHost) == 0 {
		return "localhost"
	}

	return dbHost
}

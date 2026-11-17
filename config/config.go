package config

import (
	"encoding/json"
	"os"
)

const Path = "config.dot"

type Config struct {
	Port         string `json:"port"`
	HostKeyPath  string `json:"hostKeyPath"`
	AuthKeysPath string `json:"authKeysPath"`
	Root     string `json:"root"`
  HealthSocket string `json:"healthSocket"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)

	return config, err
}

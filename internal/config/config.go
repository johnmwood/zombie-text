package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ClaudeAPIKey string `json:"claude_api_key"`
	ImageDir     string `json:"image_dir"`
	BasePrompt   string `json:"base_prompt"`
}

func LoadConfig(paths ...string) (*Config, error) {
	var cfg Config
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		if err := json.NewDecoder(file).Decode(&cfg); err != nil {
			return nil, err
		}
		file.Close()
	}
	return &cfg, nil
}

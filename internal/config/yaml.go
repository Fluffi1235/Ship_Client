package config

import (
	"diplom/internal/models"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConfigFromYaml() (*models.Config, error) {
	var cfg *models.Config

	file, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

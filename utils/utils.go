package utils

import (
	"fmt"
	"os"

	"github.com/Rokkit-exe/neewerctl/models"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}

	var config models.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return nil, err
	}

	// Use the config values as needed
	fmt.Println("Config loaded")
	return &config, nil
}

func WriteConfig(path string, config *models.Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("Error marshaling config:", err)
		return err
	}

	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return err
	}

	fmt.Println("Config saved")
	return nil
}

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

	return nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func KelvinToTemp(k int) byte {
	if k < 2900 {
		k = 2900
	}
	if k > 7000 {
		k = 7000
	}
	return byte(((k - 2900) * 40 / 4100) + 1)
}

func TempByteToKelvin(b byte) int {
	return 2900 + int(b-1)*4100/40
}

func GetProfileValues(profileName string, profiles []models.Profile) (int, int, error) {
	for _, profile := range profiles {
		if profile.Name == profileName {
			return profile.Temperature, profile.Brightness, nil
		}
	}
	return 0, 0, fmt.Errorf("profile '%s' not found", profileName)
}

package models

import "fmt"

type State struct {
	Port        string `yaml:"port"`
	Power       bool   `yaml:"power"`
	Brightness  int    `yaml:"brightness"`
	Temperature int    `yaml:"temperature"`
}

func (s State) ToString() string {
	powerStatus := "Off"
	if s.Power {
		powerStatus = "On"
	}
	return fmt.Sprintf("Power: %s\nBrightness: %d%%\nTemperature: %dK\n", powerStatus, s.Brightness, s.Temperature)
}

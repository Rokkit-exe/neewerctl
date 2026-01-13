package ctl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Port        string `json:"port"`
	Power       bool   `json:"power"`
	Brightness  int    `json:"brightness"`
	Temperature int    `json:"temperature"`
}

func defaultState() *State {
	return &State{
		Port:        "/dev/ttyUSB0",
		Power:       false,
		Brightness:  100,
		Temperature: 7000,
	}
}

func statePath() string {
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "neewerctl", "state.json")
}

func (s *State) SaveState() error {
	path := statePath()
	os.MkdirAll(filepath.Dir(path), 0o755)

	b, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(path, b, 0o644)
}

func LoadState() (*State, error) {
	b, err := os.ReadFile(statePath())
	if os.IsNotExist(err) {
		fmt.Println("State file does not exist, using default State.")
		return defaultState(), nil
	}
	if err != nil {
		return nil, err
	}

	var s State
	err = json.Unmarshal(b, &s)
	return &s, err
}

func (s State) ToString() string {
	powerStatus := "Off"
	if s.Power {
		powerStatus = "On"
	}
	return fmt.Sprintf("Power: %s\nBrightness: %d%%\nTemperature: %dK\n", powerStatus, s.Brightness, s.Temperature)
}

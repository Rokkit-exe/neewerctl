package models

type State struct {
	Port        string `yaml:"port"`
	Power       bool   `yaml:"power"`
	Brightness  int    `yaml:"brightness"`
	Temperature int    `yaml:"temperature"`
}

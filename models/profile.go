package models

type Profile struct {
	Name        string `yaml:"name"`
	Temperature int    `yaml:"temperature"`
	Brightness  int    `yaml:"brightness"`
}

package models

type Config struct {
	Devices  []Device  `yaml:"devices"`
	Profiles []Profile `yaml:"profiles"`
}

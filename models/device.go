package models

type Device struct {
	Model          string `yaml:"model"`
	VendorId       string `yaml:"vendor_id"`
	ProductId      string `yaml:"product_id"`
	Driver         string `yaml:"driver"`
	MinBrightness  int    `yaml:"min_brightness"`
	MaxBrightness  int    `yaml:"max_brightness"`
	MinTemperature int    `yaml:"min_temperature"`
	MaxTemperature int    `yaml:"max_temperature"`
	State          State  `yaml:"state"`
}

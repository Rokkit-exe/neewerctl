package main

import (
	"fmt"

	"github.com/Rokkit-exe/neewerctl/cmd"
	"github.com/Rokkit-exe/neewerctl/utils"
)

func main() {
	config, err := utils.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}
	cmd.SetConfig(config)
	cmd.Execute()
}

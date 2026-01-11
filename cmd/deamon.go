/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Rokkit-exe/neewerctl/ctl"
	"github.com/spf13/cobra"
)

// deamonCmd represents the deamon command
var deamonCmd = &cobra.Command{
	Use:   "deamon",
	Short: "Start|stop the deamon service",
	Long: `Start or stop the deamon service.
	Usage: neewerctl deamon start|stop [flags]

	Flags:
	Examples:
	- neewerctl deamon start
	- neewerctl deamon stop
`,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"start", "stop", "run", "install"},
	Run: func(cmd *cobra.Command, args []string) {
		status := args[0]
		devicePort, _ := cmd.Flags().GetString("device")

		if status == "install" {
			err := ctl.InstallService()
			if err != nil {
				fmt.Println("Error installing deamon:", err)
				return
			}
			fmt.Println("Deamon installed...")
			return
		}

		if status == "run" {
			err := ctl.RunDeamon(devicePort)
			if err != nil {
				fmt.Println("Error running deamon:", err)
				return
			}
			fmt.Println("Deamon running...")
			return
		}

		if status != "start" && status != "stop" {
			fmt.Println("Error: argument must be 'start' or 'stop'")
			cmd.Help()
			return
		}
		if status == "start" {
			err := ctl.StartDeamon()
			if err != nil {
				fmt.Println("Error starting deamon:", err)
				return
			}
			fmt.Println("Deamon started...")
			return
		}

		if status == "stop" {
			err := ctl.StopDeamon()
			if err != nil {
				fmt.Println("Error stopping deamon:", err)
				return
			}
			fmt.Println("Deamon stopped...")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(deamonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deamonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deamonCmd.Flags().StringP("device", "d", "/dev/ttyUSB0", "Specify the device port")
}

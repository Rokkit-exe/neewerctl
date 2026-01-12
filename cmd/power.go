package cmd

import (
	"fmt"
	"time"

	"github.com/Rokkit-exe/neewerctl/ctl"
	"github.com/Rokkit-exe/neewerctl/models"
	"github.com/spf13/cobra"
)

// powerCmd represents the power command
var powerCmd = &cobra.Command{
	Use:   "power [on|off]",
	Short: "Power the light on or off",
	Long: `Power on|off a device.
  *** Powering on restores the last saved brightness and temperature. ***

	Usage: neewerctl power [on|off] [flags]

	Flags:

	--device | -d [device port] 

	Examples:

	Power on the device at /dev/ttyUSB0 with saved brightness and temperature:
	- neewerctl power on --device /dev/ttyUSB0

	Power off the device at /dev/ttyUSB0:
	- neewerctl power off --device /dev/ttyUSB0
`,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"on", "off"},
	Run: func(cmd *cobra.Command, args []string) {
		state := args[0]
		devicePort, _ := cmd.Flags().GetString("device")

		var targetDevice models.Device

		// Validate the argument
		if state != "on" && state != "off" {
			fmt.Println("Error: argument must be 'on' or 'off'")
			cmd.Help()
			return
		}

		for _, dev := range Config.Devices {
			if dev.State.Port == "" {
				fmt.Println("Error: Device port not specified. Use --device flag or set in config file.")
				return
			}
			if dev.State.Port == devicePort {
				targetDevice = dev
				break
			}
		}
		deviceState, err := ctl.GetState(targetDevice.State.Port)
		if err != nil {
			fmt.Println("Error getting device state:", err)
			return
		}
		time.Sleep(500 * time.Millisecond)

		if state == "on" {
			err = ctl.Send(targetDevice.State.Port, ctl.MakeFrame(true, deviceState.Brightness, deviceState.Temperature))
			if err != nil {
				fmt.Println("Error setting saved values:", err)
				return
			}
		}

		if state == "off" {
			err := ctl.Send(targetDevice.State.Port, ctl.MakeFrame(false, deviceState.Brightness, deviceState.Temperature))
			if err != nil {
				fmt.Println("Error powering off:", err)
				return
			}
		}

		time.Sleep(500 * time.Millisecond)
		deviceState, err = ctl.GetState(targetDevice.State.Port)
		if err != nil {
			fmt.Println("Error getting device state:", err)
			return
		}
		fmt.Println(deviceState.ToString())
	},
}

func init() {
	rootCmd.AddCommand(powerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// powerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// powerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	powerCmd.Flags().StringP("device", "d", "/dev/ttyUSB0", "Serial port of the Neewer device")
}

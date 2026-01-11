package cmd

import (
	"fmt"

	"github.com/Rokkit-exe/neewerctl/ctl"
	"github.com/Rokkit-exe/neewerctl/models"
	"github.com/Rokkit-exe/neewerctl/utils"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a parameter on the device",
	Long: `Set brightness|temperature|profile.
	*** profile presets will override brightness and temperature values. ***

	Usage: neewerctl set [flags]

	Flags:

	--device | -d [device port] 
	--brightness | -b  [0-100]
	--temperature | -t [2900-7000]
	--profile | -p ["cold", "sunlight", "afternoon", "sunset", "candle"]

	Profile presets:

	- cold: Brightness 100%, Temperature 7000K
  - sunlight: Brightness 28%, Temperature 5600K
  - afternoon: Brightness 16%, Temperature 5000K
  - sunset: Brightness 16%, Temperature 4500K
  - candle: Brightness 28%, Temperature 3400K

	Examples:

	Set brightness to 80% and color temperature to 4500K to device at /dev/ttyUSB0:
	- neewerctl set --brightness 80 --temperature 4500 --device /dev/ttyUSB0

	Set predefined profile "sunset":
	- neewerctl set --profile sunset
`,
	Run: func(cmd *cobra.Command, args []string) {
		brightness, _ := cmd.Flags().GetInt("brightness")
		temperature, _ := cmd.Flags().GetInt("temperature")
		profile, _ := cmd.Flags().GetString("profile")
		devicePort, _ := cmd.Flags().GetString("device")

		var targetDevice models.Device

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

		nextTemp := targetDevice.State.Temperature
		nextBright := targetDevice.State.Brightness

		if profile != "" {
			t, b, err := ctl.GetProfileValues(profile)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			err = ctl.Send(targetDevice.State.Port, ctl.MakeFrame(true, b, t))
			if err != nil {
				fmt.Println("Error setting profile:", err)
				return
			}

			fmt.Println("Profile set to", profile)
			targetDevice.State.Temperature = t
			targetDevice.State.Brightness = b
			utils.WriteConfig("/home/frank/coding/neewerctl/config.yaml", Config)
			return
		}

		if temperature >= 0 {
			t := temperature
			if t < targetDevice.MinTemperature {
				t = targetDevice.MinTemperature
			}
			if t > targetDevice.MaxTemperature {
				t = targetDevice.MaxTemperature
			}
			nextTemp = t
		}

		if brightness >= 0 {
			b := brightness
			if b > targetDevice.MaxBrightness {
				b = targetDevice.MaxBrightness
			}
			nextBright = b
		}

		err := ctl.Send(targetDevice.State.Port, ctl.MakeFrame(true, nextBright, nextTemp))
		if err != nil {
			fmt.Println("Error setting values:", err)
			return
		}

		targetDevice.State.Temperature = nextTemp
		targetDevice.State.Brightness = nextBright
		utils.WriteConfig("/home/frank/coding/neewerctl/config.yaml", Config)
		fmt.Println("Set brightness to", nextBright, "and temperature to", nextTemp)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setCmd.Flags().IntP("brightness", "b", -1, "Set brightness level (0-100)")
	setCmd.Flags().IntP("temperature", "t", -1, "Set color temperature in Kelvin (2900-7000)")
	setCmd.Flags().StringP("profile", "p", "", "Set predefined profile (cold, sunlight, afternoon, sunset, candle)")
	setCmd.Flags().StringP("device", "d", "/dev/ttyUSB0", "Specify device port (overrides config file)")
}

package cmd

import (
	"fmt"
	"os"

	"github.com/Rokkit-exe/neewerctl/ctl"
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

		deviceState, err := ctl.LoadState()
		if err != nil {
			fmt.Println("Error loading state:", err)
			return
		}
		controller, err := ctl.NewCtl(deviceState)
		if err != nil {
			fmt.Println("Error initializing controller:", err)
			os.Exit(1)
		}
		defer controller.Close()

		state := controller.GetState()
		// Override port if specified
		if devicePort != "" {
			state.Port = devicePort
		}
		nextTemp := state.Temperature
		nextBright := state.Brightness

		if temperature >= 0 {
			nextTemp = utils.ClampInt(temperature, 2900, 7000)
		}

		if brightness >= 0 {
			nextBright = utils.ClampInt(brightness, 0, 100)
		}

		if profile != "" {
			t, b, err := utils.GetProfileValues(profile, Config.Profiles)
			if err != nil {
				fmt.Println("Error getting profile values:", err)
				return
			}
			nextTemp = t
			nextBright = b
			fmt.Println("Profile: ", profile)
		}

		err = controller.Send(true, nextBright, nextTemp)
		if err != nil {
			fmt.Println("Error setting values:", err)
			return
		}

		fmt.Println(controller.GetState().ToString())
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

package cmd

import (
	"fmt"

	"github.com/Rokkit-exe/neewerctl/ctl"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Neewer devices connected via USB serial",
	Long: `List Neewer devices connected via USB serial.

	Usage: neewerctl list [flags]

	Flags:

	--device | -d [device port] 

	Examples:

	List all connected Neewer devices:
	- neewerctl list

	List a specific device at /dev/ttyUSB0:
	- neewerctl list --device /dev/ttyUSB0

`,
	Run: func(cmd *cobra.Command, args []string) {
		devicePort, _ := cmd.Flags().GetString("device")
		if devicePort != "" {
			fmt.Println("Specified device port:", devicePort)
			fmt.Println("Not implemented: specific device listing.")
			return
		}
		ports, err := serial.GetPortsList()
		if err != nil {
			fmt.Println("Error listing serial ports:", err)
			return
		}
		if len(ports) == 0 {
			fmt.Println("No Neewer devices found.")
			return
		}

		fmt.Println("Found Neewer devices:")
		for _, port := range ports {
			state, err := ctl.GetState(port)
			if err != nil {
				continue
			}
			fmt.Println(" - Neewer PL81 Pro " + state.ToString())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().StringP("device", "d", "", "Serial port of the Neewer device")
}

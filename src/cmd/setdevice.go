package cmd

import (
	"gospt/src/tui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setDeviceCmd)
}

var setDeviceCmd = &cobra.Command{
	Use:   "setdevice",
	Short: "Shows tui to pick active device",
	Long:  `Allows setting or changing the active spotify device, shown in a tui`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartTea(ctx, client, "devices")
	},
}

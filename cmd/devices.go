package cmd

import (
	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devicesCmd)
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Prints out devices",
	Long:  `Prints out devices`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Devices(ctx, client)
	},
}

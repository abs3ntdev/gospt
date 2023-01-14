package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repeatCmd)
}

var repeatCmd = &cobra.Command{
	Use:   "repeat",
	Short: "Toggles repeat",
	Long:  `Switches between repeating your current context or not, spotifyd does not support single track loops`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Repeat(ctx, client)
	},
}

package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pauseCmd)
}

var pauseCmd = &cobra.Command{
	Use:     "pause",
	Short:   "Pauses spotify",
	Aliases: []string{"pa"},
	Long:    `Pauses currently playing song on spotify`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Pause(ctx, client)
	},
}

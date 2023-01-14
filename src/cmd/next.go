package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nextCmd)
}

var nextCmd = &cobra.Command{
	Use:     "next",
	Aliases: []string{"n"},
	Short:   "Skip to next song",
	Long:    `Skip to next song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Next(ctx, client)
	},
}

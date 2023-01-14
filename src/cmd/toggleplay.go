package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(togglePlayCmd)
}

var togglePlayCmd = &cobra.Command{
	Use:     "toggleplay",
	Aliases: []string{"t"},
	Short:   "Toggles the play state of spotify",
	Long:    `If you are playing a song it will pause and if a song is paused it will play`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.TogglePlay(ctx, client)
	},
}

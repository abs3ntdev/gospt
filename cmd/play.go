package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playCmd)
}

var playCmd = &cobra.Command{
	Use:     "play",
	Aliases: []string{"pl", "start", "s"},
	Short:   "Plays spotify",
	Long:    `Plays queued song on spotify, uses last used device and activates it if needed`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Play(ctx)
	},
}

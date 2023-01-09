package cmd

import (
	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(skipCmd)
}

var skipCmd = &cobra.Command{
	Use:   "skip",
	Short: "Skip to next song",
	Long:  `Skip to next song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Skip(ctx, client)
	},
}

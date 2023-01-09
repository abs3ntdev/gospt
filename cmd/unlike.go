package cmd

import (
	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unlikeCmd)
}

var unlikeCmd = &cobra.Command{
	Use:   "unlike",
	Short: "unlikes song",
	Long:  `unlikes song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Unlike(ctx, client)
	},
}

package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(likeCmd)
}

var likeCmd = &cobra.Command{
	Use:   "like",
	Short: "Likes song",
	Long:  `Likes song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Like(ctx, client)
	},
}

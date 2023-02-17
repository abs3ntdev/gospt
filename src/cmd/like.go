package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(likeCmd)
}

var likeCmd = &cobra.Command{
	Use:     "like",
	Aliases: []string{"l"},
	Short:   "Likes song",
	Long:    `Likes song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Like(ctx)
	},
}

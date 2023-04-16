package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Returns player status in json",
	Long:  `Returns all player status in json, useful for scripting`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Status(ctx)
	},
}

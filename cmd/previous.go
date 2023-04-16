package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(previousCmd)
}

var previousCmd = &cobra.Command{
	Use:     "previous",
	Aliases: []string{"b", "prev", "back"},
	Short:   "goes to previous song",
	Long:    `if song is playing it will start over, if close to begining of song it will go to previous song`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Previous(ctx)
	},
}

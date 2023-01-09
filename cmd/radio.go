package cmd

import (
	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(radioCmd)
}

var radioCmd = &cobra.Command{
	Use:   "radio",
	Short: "Starts radio",
	Long:  `Starts radio`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Radio(ctx, client)
	},
}

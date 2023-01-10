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
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Radio(ctx, client)
	},
}
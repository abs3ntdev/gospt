package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clearRadioCmd)
}

var clearRadioCmd = &cobra.Command{
	Use:   "clearradio",
	Short: "Wipes the radio playlist and creates an empty one",
	Long:  `Wipes the radio playlist and creates an empty one, mostly for debugging or if something goes wrong`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.ClearRadio(ctx, client)
	},
}

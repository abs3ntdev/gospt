package cmd

import (
	"gospt/internal/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refillRadioCmd)
}

var refillRadioCmd = &cobra.Command{
	Use:   "refillradio",
	Short: "Refills the radio",
	Long:  `Deletes all songs up to your position in the radio and adds that many songs to the end of the radio`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.RefillRadio(ctx, client)
	},
}

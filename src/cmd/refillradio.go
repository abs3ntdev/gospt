package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(refillRadioCmd)
}

var refillRadioCmd = &cobra.Command{
	Use:     "refillradio",
	Aliases: []string{"rr"},
	Short:   "Refills the radio",
	Long:    `Deletes all songs up to your position in the radio and adds that many songs to the end of the radio`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.RefillRadio(ctx)
	},
}

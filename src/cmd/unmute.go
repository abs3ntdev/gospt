package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unmuteCmd)
}

var unmuteCmd = &cobra.Command{
	Use:   "unmute",
	Short: "unmutes playback",
	Long:  `unmutes the spotify device, playback will continue`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := commands.SetVolume(ctx, client, 100)
		if err != nil {
			return err
		}
		return nil
	},
}

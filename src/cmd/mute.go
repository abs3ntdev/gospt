package cmd

import (
	"gospt/src/commands"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(muteCmd)
}

var muteCmd = &cobra.Command{
	Use:   "mute",
	Short: "mutes playback",
	Long:  `Mutes the spotify device, playback will continue`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := commands.SetVolume(ctx, client, 0)
		if err != nil {
			return err
		}
		return nil
	},
}

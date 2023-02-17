package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumeCmd)
}

var volumeCmd = &cobra.Command{
	Use:     "volume",
	Short:   "sets the volume",
	Aliases: []string{"v"},
	Args:    cobra.MinimumNArgs(1),
	Long:    `Sets the volume to the given percent [0-100] or increases/decreases by 5 percent if you say up or down`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "up" {
			err := commands.ChangeVolume(ctx, 5)
			if err != nil {
				return err
			}
			return nil
		}

		if args[0] == "down" {
			err := commands.ChangeVolume(ctx, -5)
			if err != nil {
				return err
			}
			return nil
		}

		vol, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		err = commands.SetVolume(ctx, vol)
		if err != nil {
			return err
		}
		return nil
	},
}

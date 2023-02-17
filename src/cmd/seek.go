package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seekCmd)
}

var seekCmd = &cobra.Command{
	Use:     "seek {forward/backward/songposition in seconds}",
	Short:   "seek forward/backward or to a given second",
	Aliases: []string{"s"},
	Args:    cobra.MinimumNArgs(1),
	Long:    `Seeks forward or backward, or seeks to a given position in seconds`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "forward" || args[0] == "f" {
			err := commands.Seek(ctx, true)
			if err != nil {
				return err
			}
			return nil
		}

		if args[0] == "backward" || args[0] == "b" {
			err := commands.Seek(ctx, false)
			if err != nil {
				return err
			}
			return nil
		}

		pos, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		pos = pos * 1000
		err = commands.SetPosition(ctx, pos)
		if err != nil {
			return err
		}
		return nil
	},
}

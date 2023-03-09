package cmd

import (
	"os"
	"path/filepath"

	"git.asdf.cafe/abs3nt/gospt/src/tui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tracksCmd)
}

var tracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "Opens saved tracks",
	Long:  `Uses TUI to open a list of saved tracks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, _ := os.UserConfigDir()
		if commands.ActiveDeviceExists(ctx) {
			return tui.StartTea(ctx, commands, "tracks")
		}
		if _, err := os.Stat(filepath.Join(configDir, "gospt/device.json")); err != nil {
			return tui.StartTea(ctx, commands, "devices")
		}
		return tui.StartTea(ctx, commands, "tracks")
	},
}

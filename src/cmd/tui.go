package cmd

import (
	"os"
	"path/filepath"

	"gitea.asdf.cafe/abs3nt/gospt/src/tui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Default command, launches the main menu",
	Long:  `Default command. this is what will run if no other commands are present. Shows the main menu.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, _ := os.UserConfigDir()
		if commands.ActiveDeviceExists(ctx) {
			return tui.StartTea(ctx, commands, "main")
		}
		if _, err := os.Stat(filepath.Join(configDir, "gospt/device.json")); err != nil {
			return tui.StartTea(ctx, commands, "devices")
		}
		return tui.StartTea(ctx, commands, "main")
	},
}

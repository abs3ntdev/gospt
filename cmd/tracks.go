package cmd

import (
	"os"
	"path/filepath"

	"gospt/internal/tui"

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
		if _, err := os.Stat(filepath.Join(configDir, "gospt/device.json")); err != nil {
			return tui.StartTea(ctx, client, "devices")
		}
		return tui.StartTea(ctx, client, "tracks")
	},
}

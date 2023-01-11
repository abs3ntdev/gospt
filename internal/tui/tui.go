package tui

import (
	"fmt"

	"gospt/internal/gctx"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(ctx *gctx.Context, client *spotify.Client, mode string) error {
	m, err := InitMain(ctx, client, mode)
	if err != nil {
		fmt.Println("UH OH")
	}
	P = tea.NewProgram(m, tea.WithAltScreen())
	if err := P.Start(); err != nil {
		return err
	}
	return nil
}

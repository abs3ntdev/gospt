package tui

import (
	"fmt"
	"log"

	"gospt/internal/gctx"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(ctx *gctx.Context, client *spotify.Client) error {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		return err
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	m, err := InitMain(ctx, client)
	if err != nil {
		fmt.Println("UH OH")
	}
	P = tea.NewProgram(m, tea.WithAltScreen())
	if err := P.Start(); err != nil {
		return err
	}
	return nil
}

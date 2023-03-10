package tui

import (
	"fmt"

	"git.asdf.cafe/abs3nt/gospt/src/commands"
	"git.asdf.cafe/abs3nt/gospt/src/gctx"

	tea "github.com/charmbracelet/bubbletea"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(ctx *gctx.Context, cmd *commands.Commands, mode string) error {
	m, err := InitMain(ctx, cmd, Mode(mode))
	if err != nil {
		fmt.Println("UH OH")
	}
	P = tea.NewProgram(m, tea.WithAltScreen())
	if err := P.Start(); err != nil {
		return err
	}
	return nil
}

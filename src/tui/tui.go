package tui

import (
	"git.asdf.cafe/abs3nt/gospt/src/commands"
	"git.asdf.cafe/abs3nt/gospt/src/gctx"

	tea "github.com/charmbracelet/bubbletea"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(ctx *gctx.Context, cmd *commands.Commands, mode string) error {
	m, err := InitMain(ctx, cmd, Mode(mode))
	if err != nil {
		return err
	}
	P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := P.Run(); err != nil {
		return err
	}
	return nil
}

package tui

import (
	"fmt"
	"os"

	"gospt/internal/commands"
	"gospt/internal/gctx"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var deviceDocStyle = lipgloss.NewStyle().Margin(1, 2)

type deviceItem struct {
	spotify.PlayerDevice
}

func (i deviceItem) Title() string { return i.Name }
func (i deviceItem) Description() string {
	return fmt.Sprintf("%s - active: %t", i.ID, i.Active)
}
func (i deviceItem) FilterValue() string { return i.Title() }

type deviceModel struct {
	list   list.Model
	page   int
	ctx    *gctx.Context
	client *spotify.Client
}

func (m deviceModel) Init() tea.Cmd {
	return nil
}

func (m deviceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			device := m.list.SelectedItem()
			var err error
			err = commands.SetDevice(m.ctx, m.client, device.(deviceItem).PlayerDevice)
			if err != nil {
				m.ctx.Printf(err.Error())
			}
			err = DisplayMain(m.ctx, m.client)
			if err != nil {
				return m, tea.Quit
			}
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if msg.Type == 5 {
			m.list.CursorUp()
		}
		if msg.Type == 6 {
			m.list.CursorDown()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m deviceModel) View() string {
	return docStyle.Render(m.list.View())
}

func DisplayDevices(ctx *gctx.Context, client *spotify.Client) error {
	items := []list.Item{}
	devices, err := client.PlayerDevices(ctx)
	if err != nil {
		return err
	}
	for _, device := range devices {
		items = append(items, deviceItem{
			device,
		})
	}
	if err != nil {
		return err
	}
	m := deviceModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0), page: 1, ctx: ctx, client: client}
	m.list.Title = "Saved Tracks"

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	fmt.Println("DEVICE SET AND SAVED")
	return nil
}

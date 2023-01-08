package tui

import (
	"fmt"
	"os"
	"time"

	"gospt/ctx"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	spotify.SavedTrack
}

func (i item) Title() string { return i.Name }
func (i item) Description() string {
	return fmt.Sprint(i.TimeDuration().Round(time.Second), " by ", i.Artists[0].Name)
}
func (i item) FilterValue() string { return i.Title() }

type model struct {
	list   list.Model
	page   int
	ctx    *ctx.Context
	client *spotify.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "ctrl+n" {
			tracks, err := TrackList(m.ctx, m.client, (m.page + 1))
			if err != nil {
				return m, tea.Quit
			}
			m.page++
			items := []list.Item{}
			for _, track := range tracks.Tracks {
				items = append(items, item{
					track,
				})
			}
			for _, item := range items {
				m.list.InsertItem(len(m.list.Items())+1, item)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func DisplayList(ctx *ctx.Context, client *spotify.Client) error {
	tracks, err := TrackList(ctx, client, 1)
	if err != nil {
		return err
	}
	items := []list.Item{}
	for _, track := range tracks.Tracks {
		items = append(items, item{
			track,
		})
	}
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), page: 1, ctx: ctx, client: client}
	m.list.Title = "Saved Tracks"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}

func TrackList(ctx *ctx.Context, client *spotify.Client, page int) (*spotify.SavedTrackPage, error) {
	return client.CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

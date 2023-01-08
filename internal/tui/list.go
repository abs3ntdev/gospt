package tui

import (
	"fmt"
	"os"
	"time"

	"gospt/internal/commands"
	"gospt/internal/gctx"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	Name     string
	Duration string
	Artist   spotify.SimpleArtist
	ID       spotify.ID
	spotify.SavedTrack
}

func (i item) Title() string { return i.Name }
func (i item) Description() string {
	return fmt.Sprint(i.Duration, " by ", i.Artist.Name)
}
func (i item) FilterValue() string { return i.Title() + i.Artist.Name }

type model struct {
	list   list.Model
	page   int
	ctx    *gctx.Context
	client *spotify.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.list.Paginator.OnLastPage() {
		tracks, err := commands.TrackList(m.ctx, m.client, (m.page + 1))
		if err != nil {
			return m, tea.Quit
		}
		m.page++
		items := []list.Item{}
		for _, track := range tracks.Tracks {
			items = append(items, item{
				Name:     track.Name,
				Artist:   track.Artists[0],
				Duration: track.TimeDuration().Round(time.Second).String(),
				ID:       track.ID,
			})
		}
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			track := m.list.SelectedItem()
			var err error
			err = m.client.QueueSong(m.ctx, track.(item).ID)
			if err != nil {
				m.ctx.Printf(err.Error())
			}
			err = m.client.Next(m.ctx)
			if err != nil {
				m.ctx.Printf(err.Error())
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

func DisplayList(ctx *gctx.Context, client *spotify.Client) error {
	items := []list.Item{}
	tracks, err := commands.TrackList(ctx, client, 1)
	if err != nil {
		return err
	}
	for _, track := range tracks.Tracks {
		items = append(items, item{
			Name:     track.Name,
			Artist:   track.Artists[0],
			Duration: track.TimeDuration().Round(time.Second).String(),
			ID:       track.ID,
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

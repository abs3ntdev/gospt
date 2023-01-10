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

var playlistsDocStyle = lipgloss.NewStyle().Margin(1, 2)

type playlistItem struct {
	Name string
	Desc string
	ID   spotify.ID
	spotify.SimplePlaylist
}

func (i playlistItem) Title() string       { return i.Name }
func (i playlistItem) Description() string { return i.Desc }
func (i playlistItem) FilterValue() string { return i.Title() + i.Desc }

type playlistModel struct {
	list   list.Model
	page   int
	ctx    *gctx.Context
	client *spotify.Client
}

func (m playlistModel) Init() tea.Cmd {
	return nil
}

func (m playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.list.Paginator.OnLastPage() {
		// if the last request was not full
		if len(m.list.Items())%50 == 0 {
			playlists, err := commands.Playlists(m.ctx, m.client, (m.page + 1))
			if err != nil {
				return m, tea.Quit
			}
			m.page++
			items := []list.Item{}
			for _, playlist := range playlists.Playlists {
				items = append(items, playlistItem{
					Name:           playlist.Name,
					Desc:           playlist.Description,
					ID:             playlist.ID,
					SimplePlaylist: playlist,
				})
			}
			for _, item := range items {
				m.list.InsertItem(len(m.list.Items())+1, item)
			}
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			os.Exit(0)
		}
		if msg.String() == "enter" {
			playlist := m.list.SelectedItem().(playlistItem).SimplePlaylist
			PlaylistTracks(m.ctx, m.client, playlist)
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

func (m playlistModel) View() string {
	return docStyle.Render(m.list.View())
}

func DisplayPlaylists(ctx *gctx.Context, client *spotify.Client) error {
	items := []list.Item{}
	playlists, err := commands.Playlists(ctx, client, 1)
	if err != nil {
		return err
	}
	for _, playlist := range playlists.Playlists {
		items = append(items, playlistItem{
			Name:           playlist.Name,
			Desc:           playlist.Description,
			ID:             playlist.ID,
			SimplePlaylist: playlist,
		})
	}
	m := playlistModel{
		list:   list.New(items, list.NewDefaultDelegate(), 0, 0),
		page:   1,
		ctx:    ctx,
		client: client,
	}
	m.list.Title = "Saved Tracks"

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}

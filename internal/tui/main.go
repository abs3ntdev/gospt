package tui

import (
	"fmt"
	"os"

	"gospt/internal/commands"
	"gospt/internal/gctx"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

type mainItem struct {
	Name        string
	Desc        string
	SpotifyItem any
}

func (i mainItem) Title() string       { return i.Name }
func (i mainItem) Description() string { return i.Desc }
func (i mainItem) FilterValue() string { return i.Title() + i.Desc }

type mainModel struct {
	list   list.Model
	page   int
	ctx    *gctx.Context
	client *spotify.Client
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				items = append(items, mainItem{
					Name:        playlist.Name,
					Desc:        playlist.Description,
					SpotifyItem: playlist,
				})
			}
			for _, item := range items {
				m.list.InsertItem(len(m.list.Items())+1, item)
			}
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "esc" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			switch m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
			case spotify.SimplePlaylist:
				playlist := m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimplePlaylist)
				p, err := InitPlaylists(m.ctx, m.client, playlist)
				if err != nil {
					return m, tea.Quit
				}
				play := tea.NewProgram(p, tea.WithAltScreen(), tea.WithMouseCellMotion())
				if _, err := play.Run(); err != nil {
					return m, tea.Quit
				}
			case *spotify.SavedTrackPage:
				p, err := InitSavedTracks(m.ctx, m.client)
				if err != nil {
					return m, tea.Quit
				}
				play := tea.NewProgram(p, tea.WithAltScreen(), tea.WithMouseCellMotion())
				if _, err := play.Run(); err != nil {
					return m, tea.Quit
				}
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

func (m mainModel) View() string {
	return docStyle.Render(m.list.View())
}

func DisplayMain(ctx *gctx.Context, client *spotify.Client) error {
	items := []list.Item{}
	saved_items, err := commands.TrackList(ctx, client, 1)
	items = append(items, mainItem{
		Name:        "Saved Tracks",
		Desc:        fmt.Sprintf("%d saved songs", saved_items.Total),
		SpotifyItem: saved_items,
	})
	playlists, err := commands.Playlists(ctx, client, 1)
	if err != nil {
		return err
	}
	for _, playlist := range playlists.Playlists {
		items = append(items, mainItem{
			Name:        playlist.Name,
			Desc:        playlist.Description,
			SpotifyItem: playlist,
		})
	}
	m := mainModel{
		list:   list.New(items, list.NewDefaultDelegate(), 0, 0),
		page:   1,
		ctx:    ctx,
		client: client,
	}
	m.list.Title = "GOSPT"

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}

func InitMain(ctx *gctx.Context, client *spotify.Client) (tea.Model, error) {
	items := []list.Item{}
	saved_items, err := commands.TrackList(ctx, client, 1)
	items = append(items, mainItem{
		Name:        "Saved Tracks",
		Desc:        fmt.Sprintf("%d saved songs", saved_items.Total),
		SpotifyItem: saved_items,
	})
	playlists, err := commands.Playlists(ctx, client, 1)
	if err != nil {
		return nil, err
	}
	for _, playlist := range playlists.Playlists {
		items = append(items, mainItem{
			Name:        playlist.Name,
			Desc:        playlist.Description,
			SpotifyItem: playlist,
		})
	}
	m := mainModel{
		list:   list.New(items, list.NewDefaultDelegate(), 0, 0),
		page:   1,
		ctx:    ctx,
		client: client,
	}
	m.list.Title = "GOSPT"
	return m, nil
}

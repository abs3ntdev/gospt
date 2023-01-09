package tui

import (
	"fmt"
	"os"
	"time"

	"gospt/internal/commands"
	"gospt/internal/gctx"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

type track struct {
	Name     string
	Duration string
	Artist   spotify.SimpleArtist
	ID       spotify.ID
	spotify.SavedTrack
}

func (i track) Title() string { return i.Name }
func (i track) Description() string {
	return fmt.Sprint(i.Duration, " by ", i.Artist.Name)
}
func (i track) FilterValue() string { return i.Title() + i.Artist.Name }

type playlistTracksModel struct {
	list     list.Model
	page     int
	ctx      *gctx.Context
	client   *spotify.Client
	playlist spotify.SimplePlaylist
}

func (m playlistTracksModel) Init() tea.Cmd {
	return nil
}

func (m playlistTracksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.list.Paginator.OnLastPage() {
		// if last request was still full request more
		if len(m.list.Items())%50 == 0 {
			tracks, err := commands.PlaylistTracks(m.ctx, m.client, m.playlist.ID, (m.page + 1))
			if err != nil {
				return m, tea.Quit
			}
			m.page++
			items := []list.Item{}
			for _, track := range tracks.Tracks {
				items = append(items, item{
					Name:     track.Track.Name,
					Artist:   track.Track.Artists[0],
					Duration: track.Track.TimeDuration().Round(time.Second).String(),
					ID:       track.Track.ID,
				})
			}
			for _, item := range items {
				m.list.InsertItem(len(m.list.Items())+1, item)
			}

		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "backspace" || msg.String() == "q" || msg.String() == "esc" {
			DisplayMain(m.ctx, m.client)
			return m, tea.Quit
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "ctrl+r" {
			track := m.list.SelectedItem()
			err := commands.RadioGivenSong(m.ctx, m.client, track.(item).ID, 0)
			if err != nil {
				return m, tea.Quit
			}
		}
		if msg.String() == "enter" {
			track := m.list.SelectedItem()
			var err error
			err = commands.QueueSong(m.ctx, m.client, track.(item).ID)
			if err != nil {
				m.ctx.Printf(err.Error())
			}
			err = commands.Skip(m.ctx, m.client)
			if err != nil {
				m.ctx.Printf(err.Error())
			}
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

func (m playlistTracksModel) View() string {
	return docStyle.Render(m.list.View())
}

func PlaylistTracks(ctx *gctx.Context, client *spotify.Client, playlist spotify.SimplePlaylist) error {
	items := []list.Item{}
	tracks, err := commands.PlaylistTracks(ctx, client, playlist.ID, 1)
	if err != nil {
		return err
	}
	for _, track := range tracks.Tracks {
		items = append(items, item{
			Name:     track.Track.Name,
			Artist:   track.Track.Artists[0],
			Duration: track.Track.TimeDuration().Round(time.Second).String(),
			ID:       track.Track.ID,
		})
	}

	m := playlistTracksModel{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		page:     1,
		ctx:      ctx,
		client:   client,
		playlist: playlist,
	}
	m.list.Title = playlist.Name

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return nil
}

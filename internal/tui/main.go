package tui

import (
	"fmt"
	"time"

	"gospt/internal/gctx"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var (
	P                *tea.Program
	DocStyle         = lipgloss.NewStyle().Margin(0, 2)
	currentlyPlaying string
	main_updates     chan *mainModel
	page             = 1
)

type mainItem struct {
	Name        string
	Duration    string
	Artist      spotify.SimpleArtist
	ID          spotify.ID
	Desc        string
	SpotifyItem any
}

func (i mainItem) Title() string       { return i.Name }
func (i mainItem) Description() string { return i.Desc }
func (i mainItem) FilterValue() string { return i.Title() + i.Desc }

type mainModel struct {
	list       list.Model
	ctx        *gctx.Context
	client     *spotify.Client
	mode       string
	playlist   spotify.SimplePlaylist
	artist     spotify.SimpleArtist
	album      spotify.SimpleAlbum
	fromArtist bool
}

func (m mainModel) Init() tea.Cmd {
	main_updates = make(chan *mainModel)
	return nil
}

func (m *mainModel) Tick() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				playing, _ := m.client.PlayerCurrentlyPlaying(m.ctx)
				if playing.Playing {
					currentlyPlaying = "Now playing " + playing.Item.Name + " by " + playing.Item.Artists[0].Name
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (m mainModel) View() string {
	return DocStyle.Render(m.list.View())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.list.NewStatusMessage(currentlyPlaying)
	select {
	case update := <-main_updates:
		m.list.SetItems(update.list.Items())
	default:
	}
	if m.list.Paginator.Page == m.list.Paginator.TotalPages-2 && m.list.Cursor() == 0 {
		// if last request was still full request more
		if len(m.list.Items())%50 == 0 {
			go m.LoadMoreItems()
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "d" {
			m.mode = "devices"
			new_items, err := DeviceView(m.ctx, m.client)
			if err != nil {
				fmt.Println(err.Error())
				return m, tea.Quit
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
			m.list.NewStatusMessage("Setting view to devices")
		}
		if msg.String() == "backspace" || msg.String() == "esc" || msg.String() == "q" {
			if m.mode == "album" {
				if m.fromArtist {
					m.mode = "albums"
					m.fromArtist = true
					m.list.NewStatusMessage("Opening " + m.artist.Name)
					new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.client)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
				} else {
					m.mode = "albums"
					m.list.NewStatusMessage("Setting view to albums")
					new_items, err := AlbumsView(m.ctx, m.client)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
				}
			} else if m.mode == "albums" {
				if m.fromArtist {
					m.mode = "artists"
					m.fromArtist = false
					m.list.NewStatusMessage("Setting view to artists")
					new_items, err := ArtistsView(m.ctx, m.client)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
				} else {
					m.mode = "main"
					m.list.NewStatusMessage("Setting view to main")
					new_items, err := MainView(m.ctx, m.client)
					if err != nil {
						fmt.Println(err.Error())
					}
					m.list.SetItems(new_items)
				}
			} else if m.mode != "main" {
				m.mode = "main"
				m.list.NewStatusMessage("Setting view to main")
				new_items, err := MainView(m.ctx, m.client)
				if err != nil {
					fmt.Println(err.Error())
				}
				m.list.SetItems(new_items)
			} else {
				return m, tea.Quit
			}
			m.list.ResetSelected()
			page = 0
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" || msg.String() == "spacebar" {
			switch m.mode {
			case "main":
				switch m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
				case *spotify.FullArtistCursorPage:
					m.mode = "artists"
					m.list.NewStatusMessage("Setting view to artists")
					new_items, err := ArtistsView(m.ctx, m.client)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
				case *spotify.SavedAlbumPage:
					m.mode = "albums"
					m.list.NewStatusMessage("Setting view to albums")
					new_items, err := AlbumsView(m.ctx, m.client)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
				case spotify.SimplePlaylist:
					m.mode = "playlist"
					playlist := m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimplePlaylist)
					m.playlist = playlist
					m.list.NewStatusMessage("Setting view to playlist " + playlist.Name)
					new_items, err := PlaylistView(m.ctx, m.client, playlist)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
				case *spotify.SavedTrackPage:
					m.mode = "tracks"
					m.list.NewStatusMessage("Setting view to saved tracks")
					new_items, err := SavedTracksView(m.ctx, m.client)
					if err != nil {
						fmt.Println(err.Error())
						return m, tea.Quit
					}
					m.list.SetItems(new_items)
					m.list.ResetSelected()
					m.list.NewStatusMessage("Setting view to tracks")
				}
			case "albums":
				m.mode = "album"
				m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
				m.list.NewStatusMessage("Opening " + m.album.Name)
				new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.client)
				if err != nil {
					fmt.Println(err.Error())
					return m, tea.Quit
				}
				m.list.SetItems(new_items)
				m.list.ResetSelected()
			case "artists":
				m.mode = "albums"
				m.fromArtist = true
				m.artist = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleArtist)
				m.list.NewStatusMessage("Opening " + m.artist.Name)
				new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.client)
				if err != nil {
					fmt.Println(err.Error())
					return m, tea.Quit
				}
				m.list.SetItems(new_items)
				m.list.ResetSelected()
			case "album":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Playing " + currentlyPlaying)
				go HandlePlay(m.ctx, m.client, &m.album.URI, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.TotalPages))
			case "playlist":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Playing " + currentlyPlaying)
				go HandlePlay(m.ctx, m.client, &m.playlist.URI, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.PerPage))
			case "tracks":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Playing " + currentlyPlaying)
				go HandlePlayLikedSong(m.ctx, m.client, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.PerPage))
			case "devices":
				go HandleSetDevice(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.PlayerDevice))
				m.list.NewStatusMessage("Setting device to " + m.list.SelectedItem().FilterValue())
				m.mode = "main"
				m.list.NewStatusMessage("Setting view to main")
				new_items, err := MainView(m.ctx, m.client)
				if err != nil {
					fmt.Println(err.Error())
				}
				m.list.SetItems(new_items)
			}
		}
		if msg.String() == "ctrl+r" {
			switch m.mode {
			case "main":
				switch m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
				case spotify.SimplePlaylist:
					currentlyPlaying = m.list.SelectedItem().FilterValue()
					m.list.NewStatusMessage("Starting radio for " + currentlyPlaying)
					go HandlePlaylistRadio(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimplePlaylist))
				case *spotify.SavedTrackPage:
					currentlyPlaying = m.list.SelectedItem().FilterValue()
					m.list.NewStatusMessage("Starting radio for " + currentlyPlaying)
					go HandleLibraryRadio(m.ctx, m.client)
				}
			case "albums":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Stating radio for" + currentlyPlaying)
				go HandleAlbumRadio(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum).ID)
			case "album":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Stating radio for" + currentlyPlaying)
				go HandleRadio(m.ctx, m.client, m.list.SelectedItem().(mainItem).ID)
			case "playlist":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Starting radio for " + currentlyPlaying)
				go HandleRadio(m.ctx, m.client, m.list.SelectedItem().(mainItem).ID)
			case "tracks":
				currentlyPlaying = m.list.SelectedItem().FilterValue()
				m.list.NewStatusMessage("Playing " + currentlyPlaying)
				go HandleRadio(m.ctx, m.client, m.list.SelectedItem().(mainItem).ID)
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
		h, v := DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func InitMain(ctx *gctx.Context, client *spotify.Client, mode string) (tea.Model, error) {
	playing, _ := client.PlayerCurrentlyPlaying(ctx)
	if playing.Playing {
		currentlyPlaying = "Now playing " + playing.Item.Name + " by " + playing.Item.Artists[0].Name
	}
	items := []list.Item{}
	var err error
	switch mode {
	case "main":
		items, err = MainView(ctx, client)
		if err != nil {
			return nil, err
		}
	case "devices":
		items, err = DeviceView(ctx, client)
		if err != nil {
			return nil, err
		}
	case "tracks":
		items, err = SavedTracksView(ctx, client)
		if err != nil {
			return nil, err
		}
	}
	m := mainModel{
		list:   list.New(items, list.NewDefaultDelegate(), 0, 0),
		ctx:    ctx,
		client: client,
		mode:   mode,
	}
	m.list.Title = "GOSPT"
	go m.Tick()
	m.list.DisableQuitKeybindings()
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(key.WithKeys("ctrl"+"r"), key.WithHelp("ctrl+r", "start radio")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "select device")),
		}
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(key.WithKeys("ctrl"+"r"), key.WithHelp("ctrl+r", "start radio")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "select device")),
		}
	}
	return m, nil
}

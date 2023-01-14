package tui

import (
	"fmt"
	"time"

	"gospt/src/gctx"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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

type Mode string

const (
	Album             Mode = "album"
	ArtistAlbum            = "artistalbum"
	Artist                 = "artist"
	Artists                = "artists"
	Tracks                 = "tracks"
	Albums                 = "albums"
	Main                   = "main"
	Playlists              = "playlists"
	Playlist               = "playlist"
	Devices                = "devices"
	Search                 = "search"
	SearchAlbums           = "searchalbums"
	SearchAlbum            = "searchalbum"
	SearchArtists          = "searchartists"
	SearchArtist           = "searchartist"
	SearchArtistAlbum      = "searchartistalbum"
	SearchTracks           = "searchtracks"
	SearchPlaylists        = "searchplaylsits"
	SearchPlaylist         = "searchplaylist"
)

type mainItem struct {
	Name        string
	Duration    string
	Artist      spotify.SimpleArtist
	ID          spotify.ID
	Desc        string
	SpotifyItem any
}

type SearchResults struct {
	Tracks    *spotify.FullTrackPage
	Artists   *spotify.FullArtistPage
	Playlists *spotify.SimplePlaylistPage
	Albums    *spotify.SimpleAlbumPage
}

func (i mainItem) Title() string       { return i.Name }
func (i mainItem) Description() string { return i.Desc }
func (i mainItem) FilterValue() string { return i.Title() + i.Desc }

type mainModel struct {
	list          list.Model
	input         textinput.Model
	ctx           *gctx.Context
	client        *spotify.Client
	mode          Mode
	playlist      spotify.SimplePlaylist
	artist        spotify.SimpleArtist
	album         spotify.SimpleAlbum
	searchResults *SearchResults
	search        string
}

func (m *mainModel) PlayRadio() {
	currentlyPlaying = m.list.SelectedItem().FilterValue()
	m.list.NewStatusMessage("Playing " + currentlyPlaying)
	selectedItem := m.list.SelectedItem().(mainItem).SpotifyItem
	switch selectedItem.(type) {
	case spotify.SimplePlaylist:
		go HandlePlaylistRadio(m.ctx, m.client, selectedItem.(spotify.SimplePlaylist))
		return
	case spotify.SavedTrackPage:
		go HandleLibraryRadio(m.ctx, m.client)
		return
	case spotify.SimpleAlbum:
		go HandleAlbumRadio(m.ctx, m.client, selectedItem.(spotify.SimpleAlbum).ID)
		return
	case spotify.FullAlbum:
		go HandleAlbumRadio(m.ctx, m.client, selectedItem.(spotify.FullAlbum).ID)
		return
	case spotify.SimpleArtist:
		go HandleArtistRadio(m.ctx, m.client, selectedItem.(spotify.SimpleArtist).ID)
		return
	case spotify.FullArtist:
		go HandleArtistRadio(m.ctx, m.client, selectedItem.(spotify.FullArtist).ID)
		return
	default:
		go HandleRadio(m.ctx, m.client, m.list.SelectedItem().(mainItem).ID)
		return
	}
}

func (m *mainModel) GoBack() (tea.Cmd, error) {
	switch m.mode {
	case Main:
		return tea.Quit, nil
	case Albums, Artists, Tracks, Playlist, Devices, Search:
		m.mode = Main
		m.list.NewStatusMessage("Setting view to main")
		new_items, err := MainView(m.ctx, m.client)
		if err != nil {
			fmt.Println(err.Error())
		}
		m.list.SetItems(new_items)

	case Album:
		m.mode = Albums
		m.list.NewStatusMessage("Setting view to albums")
		new_items, err := AlbumsView(m.ctx, m.client)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()

	case Artist:
		m.mode = Artists
		m.list.NewStatusMessage("Setting view to artists")
		new_items, err := ArtistsView(m.ctx, m.client)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()

	case ArtistAlbum:
		m.mode = Artist
		m.list.NewStatusMessage("Opening " + m.artist.Name)
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.client)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()

	case SearchArtists, SearchTracks, SearchAlbums, SearchPlaylists:
		m.mode = Search
		m.list.NewStatusMessage("Setting view to search for " + m.input.Value())
		items, result, err := SearchView(m.ctx, m.client, m.search)
		if err != nil {
			return nil, err
		}
		m.searchResults = result
		m.list.SetItems(items)
	case SearchArtist:
		m.mode = SearchArtists
		m.list.NewStatusMessage("Setting view to artists")
		new_items, err := SearchArtistsView(m.ctx, m.client, m.searchResults.Artists)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case SearchArtistAlbum:
		m.mode = SearchArtist
		m.list.NewStatusMessage("Opening " + m.artist.Name)
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.client)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchAlbum:
		m.mode = SearchAlbums
		m.list.NewStatusMessage("Setting view to albums")
		new_items, err := SearchAlbumsView(m.ctx, m.client, m.searchResults.Albums)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case SearchPlaylist:
		m.mode = SearchPlaylists
		m.list.NewStatusMessage("Setting view to playlists")
		new_items, err := SearchPlaylistsView(m.ctx, m.client, m.searchResults.Playlists)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	default:
		m.list.ResetSelected()
		page = 0
	}
	return nil, nil
}

func (m *mainModel) SelectItem() error {
	switch m.mode {
	case Search:
		switch m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
		case *spotify.FullArtistPage:
			m.mode = SearchArtists
			m.list.NewStatusMessage("Setting view to artists")
			new_items, err := SearchArtistsView(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(*spotify.FullArtistPage))
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SimpleAlbumPage:
			m.mode = SearchAlbums
			m.list.NewStatusMessage("Setting view to albums")
			new_items, err := SearchAlbumsView(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(*spotify.SimpleAlbumPage))
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SimplePlaylistPage:
			m.mode = SearchPlaylists
			playlists := m.list.SelectedItem().(mainItem).SpotifyItem.(*spotify.SimplePlaylistPage)
			m.list.NewStatusMessage("Setting view to playlist")
			new_items, err := SearchPlaylistsView(m.ctx, m.client, playlists)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.FullTrackPage:
			m.mode = SearchTracks
			m.list.NewStatusMessage("Setting view to tracks")
			new_items, err := SearchTracksView(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(*spotify.FullTrackPage))
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
			m.list.NewStatusMessage("Setting view to tracks")
		}
	case SearchArtists:
		m.mode = SearchArtist
		m.artist = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleArtist)
		m.list.NewStatusMessage("Opening " + m.artist.Name)
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.client)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchArtist:
		m.mode = SearchArtistAlbum
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		m.list.NewStatusMessage("Opening " + m.album.Name)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.client)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchAlbums:
		m.mode = SearchAlbum
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		m.list.NewStatusMessage("Opening " + m.album.Name)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.client)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchPlaylists:
		m.mode = SearchPlaylist
		playlist := m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimplePlaylist)
		m.playlist = playlist
		m.list.NewStatusMessage("Setting view to playlist " + playlist.Name)
		new_items, err := PlaylistView(m.ctx, m.client, playlist)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Main:
		switch m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
		case *spotify.FullArtistCursorPage:
			m.mode = Artists
			m.list.NewStatusMessage("Setting view to artists")
			new_items, err := ArtistsView(m.ctx, m.client)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SavedAlbumPage:
			m.mode = Albums
			m.list.NewStatusMessage("Setting view to albums")
			new_items, err := AlbumsView(m.ctx, m.client)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case spotify.SimplePlaylist:
			m.mode = Playlist
			playlist := m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimplePlaylist)
			m.playlist = playlist
			m.list.NewStatusMessage("Setting view to playlist " + playlist.Name)
			new_items, err := PlaylistView(m.ctx, m.client, playlist)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SavedTrackPage:
			m.mode = Tracks
			m.list.NewStatusMessage("Setting view to saved tracks")
			new_items, err := SavedTracksView(m.ctx, m.client)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
			m.list.NewStatusMessage("Setting view to tracks")
		}
	case Albums:
		m.mode = Album
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		m.list.NewStatusMessage("Opening " + m.album.Name)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.client)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Artist:
		m.mode = ArtistAlbum
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		m.list.NewStatusMessage("Opening " + m.album.Name)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.client)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Artists:
		m.mode = Artist
		m.artist = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleArtist)
		m.list.NewStatusMessage("Opening " + m.artist.Name)
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.client)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Album, ArtistAlbum, SearchArtistAlbum, SearchAlbum:
		currentlyPlaying = m.list.SelectedItem().FilterValue()
		m.list.NewStatusMessage("Playing " + currentlyPlaying)
		go HandlePlayWithContext(m.ctx, m.client, &m.album.URI, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.TotalPages))
	case Playlist, SearchPlaylist:
		currentlyPlaying = m.list.SelectedItem().FilterValue()
		m.list.NewStatusMessage("Playing " + currentlyPlaying)
		go HandlePlayWithContext(m.ctx, m.client, &m.playlist.URI, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.PerPage))
	case Tracks:
		currentlyPlaying = m.list.SelectedItem().FilterValue()
		m.list.NewStatusMessage("Playing " + currentlyPlaying)
		go HandlePlayLikedSong(m.ctx, m.client, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.PerPage))
	case SearchTracks:
		currentlyPlaying = m.list.SelectedItem().FilterValue()
		m.list.NewStatusMessage("Playing " + currentlyPlaying)
		go HandlePlayTrack(m.ctx, m.client, m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.FullTrack).ID)
	case Devices:
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
	return nil
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
	if m.input.Focused() {
		return DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return DocStyle.Render(m.list.View() + "\n")
}

func (m *mainModel) Typing(msg tea.KeyMsg) (bool, tea.Cmd) {
	if msg.String() == "enter" {
		m.list.NewStatusMessage("Setting view to search for " + m.input.Value())
		items, result, err := SearchView(m.ctx, m.client, m.input.Value())
		if err != nil {
			fmt.Println(err.Error())
			return false, tea.Quit
		}
		m.searchResults = result
		m.search = m.input.Value()
		m.list.SetItems(items)
		m.list.ResetSelected()
		m.input.SetValue("")
		m.input.Blur()
		return true, nil
	}
	if msg.String() == "esc" {
		m.input.SetValue("")
		m.input.Blur()
		return false, nil
	}
	m.input, _ = m.input.Update(msg)
	return false, nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update Now Playing
	m.list.NewStatusMessage(currentlyPlaying)
	// Update list items from LoadMore
	select {
	case update := <-main_updates:
		m.list.SetItems(update.list.Items())
	default:
	}
	// Call for more items if needed
	if m.list.Paginator.Page == m.list.Paginator.TotalPages-2 && m.list.Cursor() == 0 {
		// if last request was still full request more
		if len(m.list.Items())%50 == 0 {
			go m.LoadMoreItems()
		}
	}
	// Handle user input
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// search input
		if m.input.Focused() {
			search, cmd := m.Typing(msg)
			if search {
				m.mode = "search"
			}
			return m, cmd
		}
		// start search
		if msg.String() == "s" || msg.String() == "/" {
			m.input.Focus()
		}
		// enter device selection
		if msg.String() == "d" {
			m.mode = Devices
			new_items, err := DeviceView(m.ctx, m.client)
			if err != nil {
				fmt.Println(err.Error())
				return m, tea.Quit
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
			m.list.NewStatusMessage("Setting view to devices")
		}
		// go back
		if msg.String() == "backspace" || msg.String() == "esc" || msg.String() == "q" {
			msg, err := m.GoBack()
			if err != nil {
				fmt.Println(err)
			}
			return m, msg
		}

		// select item
		if msg.String() == "enter" || msg.String() == "spacebar" {
			err := m.SelectItem()
			if err != nil {
				return m, tea.Quit
			}
		}
		// start radio
		if msg.String() == "ctrl+r" {
			m.PlayRadio()
		}

	// handle mouse
	case tea.MouseMsg:
		if msg.Type == 5 {
			m.list.CursorUp()
		}
		if msg.Type == 6 {
			m.list.CursorDown()
		}

	// window size -1 to handle search bar
	case tea.WindowSizeMsg:
		h, v := DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-1)
	}

	// return
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func InitMain(ctx *gctx.Context, client *spotify.Client, mode Mode) (tea.Model, error) {
	playing, _ := client.PlayerCurrentlyPlaying(ctx)
	if playing.Playing {
		currentlyPlaying = "Now playing " + playing.Item.Name + " by " + playing.Item.Artists[0].Name
	}
	items := []list.Item{}
	var err error
	switch mode {
	case Main:
		items, err = MainView(ctx, client)
		if err != nil {
			return nil, err
		}
	case Devices:
		items, err = DeviceView(ctx, client)
		if err != nil {
			return nil, err
		}
	case Tracks:
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
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "search")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(key.WithKeys("ctrl"+"r"), key.WithHelp("ctrl+r", "start radio")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "select device")),
		}
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "search")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(key.WithKeys("ctrl"+"r"), key.WithHelp("ctrl+r", "start radio")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "select device")),
		}
	}
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Search..."
	input.CharLimit = 250
	input.Width = 50
	m.input = input
	return m, nil
}
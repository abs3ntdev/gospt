package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"

	"git.asdf.cafe/abs3nt/gospt/src/commands"
	"git.asdf.cafe/abs3nt/gospt/src/gctx"
)

var (
	P                *tea.Program
	DocStyle         = lipgloss.NewStyle().Margin(0, 2).Border(lipgloss.DoubleBorder(), true, true, true, true)
	currentlyPlaying *spotify.CurrentlyPlaying
	playbackContext  string
	main_updates     chan *mainModel
	page             = 1
	loading          = false
	showingMessage   = false
)

type Mode string

const (
	Album             Mode = "album"
	ArtistAlbum       Mode = "artistalbum"
	Artist            Mode = "artist"
	Artists           Mode = "artists"
	Queue             Mode = "queue"
	Tracks            Mode = "tracks"
	Albums            Mode = "albums"
	Main              Mode = "main"
	Playlists         Mode = "playlists"
	Playlist          Mode = "playlist"
	Devices           Mode = "devices"
	Search            Mode = "search"
	SearchAlbums      Mode = "searchalbums"
	SearchAlbum       Mode = "searchalbum"
	SearchArtists     Mode = "searchartists"
	SearchArtist      Mode = "searchartist"
	SearchArtistAlbum Mode = "searchartistalbum"
	SearchTracks      Mode = "searchtracks"
	SearchPlaylists   Mode = "searchplaylsits"
	SearchPlaylist    Mode = "searchplaylist"
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
	list            list.Model
	input           textinput.Model
	ctx             *gctx.Context
	commands        *commands.Commands
	mode            Mode
	playlist        spotify.SimplePlaylist
	artist          spotify.SimpleArtist
	album           spotify.SimpleAlbum
	searchResults   *SearchResults
	progress        progress.Model
	playing         *spotify.CurrentlyPlaying
	playbackContext string
	search          string
}

func (m *mainModel) PlayRadio() {
	go m.SendMessage("Starting radio for "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
	selectedItem := m.list.SelectedItem().(mainItem).SpotifyItem
	switch item := selectedItem.(type) {
	case spotify.SimplePlaylist:
		go HandlePlaylistRadio(m.ctx, m.commands, item)
		return
	case *spotify.SavedTrackPage:
		go HandleLibraryRadio(m.ctx, m.commands)
		return
	case spotify.SimpleAlbum:
		go HandleAlbumRadio(m.ctx, m.commands, item)
		return
	case spotify.FullAlbum:
		go HandleAlbumRadio(m.ctx, m.commands, item.SimpleAlbum)
		return
	case spotify.SimpleArtist:
		go HandleArtistRadio(m.ctx, m.commands, item)
		return
	case spotify.FullArtist:
		go HandleArtistRadio(m.ctx, m.commands, item.SimpleArtist)
		return
	case spotify.SimpleTrack:
		go HandleRadio(m.ctx, m.commands, item)
		return
	case spotify.FullTrack:
		go HandleRadio(m.ctx, m.commands, item.SimpleTrack)
		return
	case spotify.PlaylistTrack:
		go HandleRadio(m.ctx, m.commands, item.Track.SimpleTrack)
		return
	case spotify.PlaylistItem:
		go HandleRadio(m.ctx, m.commands, item.Track.Track.SimpleTrack)
		return
	case spotify.SavedTrack:
		go HandleRadio(m.ctx, m.commands, item.SimpleTrack)
		return
	}
}

func (m *mainModel) GoBack() (tea.Cmd, error) {
	page = 1
	switch m.mode {
	case Main:
		return tea.Quit, nil
	case Albums, Artists, Tracks, Playlist, Devices, Search, Queue:
		m.mode = Main
		new_items, err := MainView(m.ctx, m.commands)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case Album:
		m.mode = Albums
		new_items, err := AlbumsView(m.ctx, m.commands)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case Artist:
		m.mode = Artists
		new_items, err := ArtistsView(m.ctx, m.commands)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case ArtistAlbum:
		m.mode = Artist
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.commands)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case SearchArtists, SearchTracks, SearchAlbums, SearchPlaylists:
		m.mode = Search
		items, result, err := SearchView(m.ctx, m.commands, m.search)
		if err != nil {
			return nil, err
		}
		m.searchResults = result
		m.list.SetItems(items)
	case SearchArtist:
		m.mode = SearchArtists
		new_items, err := SearchArtistsView(m.ctx, m.commands, m.searchResults.Artists)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case SearchArtistAlbum:
		m.mode = SearchArtist
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.commands)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case SearchAlbum:
		m.mode = SearchAlbums
		new_items, err := SearchAlbumsView(m.ctx, m.commands, m.searchResults.Albums)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	case SearchPlaylist:
		m.mode = SearchPlaylists
		new_items, err := SearchPlaylistsView(m.ctx, m.commands, m.searchResults.Playlists)
		if err != nil {
			return nil, err
		}
		m.list.SetItems(new_items)
	default:
		page = 0
	}
	return nil, nil
}

type SpotifyUrl struct {
	ExternalURLs map[string]string
}

func (m *mainModel) CopyToClipboard() error {
	item := m.list.SelectedItem().(mainItem).SpotifyItem
	switch converted := item.(type) {
	case spotify.SimplePlaylist:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case *spotify.FullPlaylist:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case spotify.SimpleAlbum:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case *spotify.FullAlbum:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case spotify.SimpleArtist:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case *spotify.FullArtist:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case spotify.SimpleTrack:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case spotify.PlaylistTrack:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.Track.ExternalURLs["spotify"])
	case spotify.SavedTrack:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	case spotify.FullTrack:
		go m.SendMessage("Copying link to "+m.list.SelectedItem().(mainItem).Title(), 2*time.Second)
		clipboard.WriteAll(converted.ExternalURLs["spotify"])
	}
	return nil
}

func (m *mainModel) SendMessage(msg string, duration time.Duration) {
	showingMessage = true
	defer func() {
		showingMessage = false
	}()
	m.list.NewStatusMessage(msg)
	time.Sleep(duration)
}

func (m *mainModel) QueueItem() error {
	switch item := m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
	case spotify.PlaylistTrack:
		go m.SendMessage("Adding "+item.Track.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.Track.ID)
	case spotify.SavedTrack:
		go m.SendMessage("Adding "+item.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.ID)
	case spotify.SimpleTrack:
		go m.SendMessage("Adding "+item.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.ID)
	case spotify.FullTrack:
		go m.SendMessage("Adding "+item.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.ID)
	case *spotify.FullTrack:
		go m.SendMessage("Adding "+item.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.ID)
	case *spotify.SimpleTrack:
		go m.SendMessage("Adding "+item.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.ID)
	case *spotify.SimplePlaylist:
		go m.SendMessage("Adding "+item.Name+" to queue", 2*time.Second)
		go HandleQueueItem(m.ctx, m.commands, item.ID)
	}
	if m.mode == Queue {
		go func() {
			new_items, err := QueueView(m.ctx, m.commands)
			if err != nil {
				return
			}
			m.list.SetItems(new_items)
		}()
	}
	return nil
}

func (m *mainModel) DeleteTrackFromPlaylist() error {
	if m.mode != Playlist {
		return nil
	}
	track := m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.PlaylistTrack).Track
	go m.SendMessage("Deleteing "+track.Name+" from "+m.playlist.Name, 2*time.Second)
	go func() {
		HandleDeleteTrackFromPlaylist(m.ctx, m.commands, track.ID, m.playlist.ID)
		new_items, err := PlaylistView(m.ctx, m.commands, m.playlist)
		if err != nil {
			return
		}
		m.list.SetItems(new_items)
	}()
	return nil
}

func (m *mainModel) SelectItem() error {
	switch m.mode {
	case Queue:
		page = 1
		go func() {
			HandleNextInQueue(m.ctx, m.commands, m.list.Index())
			new_items, err := QueueView(m.ctx, m.commands)
			if err != nil {
				return
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		}()
	case Search:
		page = 1
		switch item := m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
		case *spotify.FullArtistPage:
			m.mode = SearchArtists
			new_items, err := SearchArtistsView(m.ctx, m.commands, item)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SimpleAlbumPage:
			m.mode = SearchAlbums
			new_items, err := SearchAlbumsView(m.ctx, m.commands, item)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SimplePlaylistPage:
			m.mode = SearchPlaylists
			new_items, err := SearchPlaylistsView(m.ctx, m.commands, item)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.FullTrackPage:
			m.mode = SearchTracks
			new_items, err := SearchTracksView(m.ctx, m.commands, item)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		}
	case SearchArtists:
		page = 1
		m.mode = SearchArtist
		m.artist = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleArtist)
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchArtist:
		page = 1
		m.mode = SearchArtistAlbum
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchAlbums:
		page = 1
		m.mode = SearchAlbum
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case SearchPlaylists:
		page = 1
		m.mode = SearchPlaylist
		playlist := m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimplePlaylist)
		m.playlist = playlist
		new_items, err := PlaylistView(m.ctx, m.commands, playlist)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Main:
		page = 1
		switch item := m.list.SelectedItem().(mainItem).SpotifyItem.(type) {
		case spotify.Queue:
			m.mode = Queue
			new_items, err := QueueView(m.ctx, m.commands)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.FullArtistCursorPage:
			m.mode = Artists
			new_items, err := ArtistsView(m.ctx, m.commands)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SavedAlbumPage:
			m.mode = Albums
			new_items, err := AlbumsView(m.ctx, m.commands)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case spotify.SimplePlaylist:
			m.mode = Playlist
			m.playlist = item
			new_items, err := PlaylistView(m.ctx, m.commands, item)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		case *spotify.SavedTrackPage:
			m.mode = Tracks
			new_items, err := SavedTracksView(m.ctx, m.commands)
			if err != nil {
				return err
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		}
	case Albums:
		page = 1
		m.mode = Album
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Artist:
		m.mode = ArtistAlbum
		m.album = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleAlbum)
		new_items, err := AlbumTracksView(m.ctx, m.album.ID, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Artists:
		m.mode = Artist
		m.artist = m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.SimpleArtist)
		new_items, err := ArtistAlbumsView(m.ctx, m.artist.ID, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
		m.list.ResetSelected()
	case Album, ArtistAlbum, SearchArtistAlbum, SearchAlbum:
		pos := m.list.Cursor() + (m.list.Paginator.Page * m.list.Paginator.TotalPages)
		go HandlePlayWithContext(m.ctx, m.commands, &m.album.URI, &pos)
	case Playlist, SearchPlaylist:
		pos := m.list.Cursor() + (m.list.Paginator.Page * m.list.Paginator.PerPage)
		go HandlePlayWithContext(m.ctx, m.commands, &m.playlist.URI, &pos)
	case Tracks:
		go HandlePlayLikedSong(m.ctx, m.commands, m.list.Cursor()+(m.list.Paginator.Page*m.list.Paginator.PerPage))
	case SearchTracks:
		go HandlePlayTrack(m.ctx, m.commands, m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.FullTrack).ID)
	case Devices:
		go HandleSetDevice(m.ctx, m.commands, m.list.SelectedItem().(mainItem).SpotifyItem.(spotify.PlayerDevice))
		go m.SendMessage("Setting device to "+m.list.SelectedItem().FilterValue(), 2*time.Second)
		m.mode = "main"
		new_items, err := MainView(m.ctx, m.commands)
		if err != nil {
			return err
		}
		m.list.SetItems(new_items)
	}
	return nil
}

func (m *mainModel) Init() tea.Cmd {
	main_updates = make(chan *mainModel)
	return Tick()
}

type tickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *mainModel) TickPlayback() {
	playing, _ := m.commands.Client().PlayerCurrentlyPlaying(m.ctx)
	if playing != nil && playing.Playing && playing.Item != nil {
		currentlyPlaying = playing
		playbackContext, _ = m.getContext(playing)
	}
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				playing, _ := m.commands.Client().PlayerCurrentlyPlaying(m.ctx)
				if playing != nil && playing.Playing && playing.Item != nil {
					currentlyPlaying = playing
					playbackContext, _ = m.getContext(playing)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *mainModel) View() string {
	if m.input.Focused() {
		return DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return DocStyle.Render(m.list.View() + "\n")
}

func (m *mainModel) Typing(msg tea.KeyMsg) (bool, tea.Cmd) {
	if msg.String() == "enter" {
		items, result, err := SearchView(m.ctx, m.commands, m.input.Value())
		if err != nil {
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

func (m *mainModel) getContext(playing *spotify.CurrentlyPlaying) (string, error) {
	context := playing.PlaybackContext
	uri_split := strings.Split(string(context.URI), ":")
	if len(uri_split) < 3 {
		return "", fmt.Errorf("NO URI")
	}
	id := strings.Split(string(context.URI), ":")[2]
	switch context.Type {
	case "album":
		album, err := m.commands.Client().GetAlbum(m.ctx, spotify.ID(id))
		if err != nil {
			return "", err
		}
		return album.Name, nil
	case "playlist":
		playlist, err := m.commands.Client().GetPlaylist(m.ctx, spotify.ID(id))
		if err != nil {
			return "", err
		}
		return playlist.Name, nil
	case "artist":
		artist, err := m.commands.Client().GetArtist(m.ctx, spotify.ID(id))
		if err != nil {
			return "", err
		}
		return artist.Name, nil
	}
	return "", nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Update list items from LoadMore
	select {
	case update := <-main_updates:
		m.list.SetItems(update.list.Items())
	default:
	}
	// Call for more items if needed
	if m.list.Paginator.Page == m.list.Paginator.TotalPages-1 && m.list.Cursor() == 0 && !loading {
		// if last request was still full request more
		if len(m.list.Items())%50 == 0 {
			go m.LoadMoreItems()
		}
	}
	// Handle user input
	switch msg := msg.(type) {
	case tickMsg:
		playing := currentlyPlaying
		if playing != nil && playing.Playing && playing.Item != nil {
			cmd := m.progress.SetPercent(float64(playing.Progress) / float64(playing.Item.Duration))
			m.playing = playing
			m.playbackContext = playbackContext
			if m.mode == Queue && len(m.list.Items()) != 0 {
				if m.list.Items()[0].(mainItem).SpotifyItem.(spotify.FullTrack).Name != playing.Item.Name {
					go func() {
						new_items, err := QueueView(m.ctx, m.commands)
						if err != nil {
							return
						}
						m.list.SetItems(new_items)
					}()
				}
			}
			return m, tea.Batch(Tick(), cmd)
		}
		return m, Tick()

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		if !showingMessage {
			m.list.NewStatusMessage(
				fmt.Sprintf("Now playing %s by %s - %s %s/%s : %s",
					m.playing.Item.Name,
					m.playing.Item.Artists[0].Name,
					m.progress.View(),
					(time.Duration(m.playing.Progress) * time.Millisecond).Round(time.Second),
					(time.Duration(m.playing.Item.Duration) * time.Millisecond).Round(time.Second),
					m.playbackContext),
			)
		}
		return m, cmd
	case tea.KeyMsg:
		// quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "c" {
			err := m.CopyToClipboard()
			if err != nil {
				return m, tea.Quit
			}
		}
		if msg.String() == ">" {
			go HandleSeek(m.ctx, m.commands, true)
		}
		if msg.String() == "<" {
			go HandleSeek(m.ctx, m.commands, false)
		}
		if msg.String() == "+" {
			go HandleVolume(m.ctx, m.commands, true)
		}
		if msg.String() == "-" {
			go HandleVolume(m.ctx, m.commands, false)
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
			new_items, err := DeviceView(m.ctx, m.commands)
			if err != nil {
				return m, tea.Quit
			}
			m.list.SetItems(new_items)
			m.list.ResetSelected()
		}
		// go back
		if msg.String() == "backspace" || msg.String() == "esc" || msg.String() == "q" {
			msg, err := m.GoBack()
			if err != nil {
				return m, tea.Quit
			}
			m.list.ResetSelected()
			return m, msg
		}
		if msg.String() == "ctrl+d" {
			err := m.DeleteTrackFromPlaylist()
			if err != nil {
				return m, tea.Quit
			}
		}
		if msg.String() == "ctrl+@" || msg.String() == "ctrl+p" {
			err := m.QueueItem()
			if err != nil {
				return m, tea.Quit
			}
		}
		// select item
		if msg.String() == "enter" || msg.String() == " " || msg.String() == "p" {
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

func InitMain(ctx *gctx.Context, c *commands.Commands, mode Mode) (tea.Model, error) {
	prog := progress.New(progress.WithColorProfile(2), progress.WithoutPercentage())
	var err error
	lipgloss.SetColorProfile(2)
	items := []list.Item{}
	switch mode {
	case Main:
		items, err = MainView(ctx, c)
		if err != nil {
			return nil, err
		}
	case Devices:
		items, err = DeviceView(ctx, c)
		if err != nil {
			return nil, err
		}
	case Tracks:
		items, err = SavedTracksView(ctx, c)
		if err != nil {
			return nil, err
		}
	}
	m := &mainModel{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		ctx:      ctx,
		commands: c,
		mode:     mode,
		progress: prog,
	}
	m.list.Title = "GOSPT"
	go m.TickPlayback()
	Tick()
	m.list.DisableQuitKeybindings()
	m.list.SetFilteringEnabled(false)
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
			key.NewBinding(key.WithKeys("ctrl"+"r"), key.WithHelp("ctrl+r", "radio")),
			key.NewBinding(key.WithKeys("ctrl"+"p"), key.WithHelp("ctrl+p", "queue")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "select device")),
		}
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
			key.NewBinding(key.WithKeys(">"), key.WithHelp(">", "seek forward")),
			key.NewBinding(key.WithKeys("<"), key.WithHelp("<", "seek backward")),
			key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "volume up")),
			key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "volume down")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy link to item")),
			key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			key.NewBinding(key.WithKeys("ctrl"+"r"), key.WithHelp("ctrl+r", "start radio")),
			key.NewBinding(key.WithKeys("ctrl"+"p"), key.WithHelp("ctrl+p", "queue song")),
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

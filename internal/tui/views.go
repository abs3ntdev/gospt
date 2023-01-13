package tui

import (
	"fmt"
	"sync"
	"time"

	"gospt/internal/commands"
	"gospt/internal/gctx"

	"github.com/charmbracelet/bubbles/list"
	"github.com/zmb3/spotify/v2"
)

func DeviceView(ctx *gctx.Context, client *spotify.Client) ([]list.Item, error) {
	items := []list.Item{}
	devices, err := client.PlayerDevices(ctx)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		items = append(items, mainItem{
			Name:        device.Name,
			Desc:        fmt.Sprintf("%s - active: %t", device.ID, device.Active),
			SpotifyItem: device,
		})
	}
	return items, nil
}

func PlaylistView(ctx *gctx.Context, client *spotify.Client, playlist spotify.SimplePlaylist) ([]list.Item, error) {
	items := []list.Item{}
	tracks, err := commands.PlaylistTracks(ctx, client, playlist.ID, 1)
	if err != nil {
		return nil, err
	}
	for _, track := range tracks.Tracks {
		items = append(items, mainItem{
			Name:     track.Track.Name,
			Artist:   track.Track.Artists[0],
			Duration: track.Track.TimeDuration().Round(time.Second).String(),
			ID:       track.Track.ID,
			Desc:     track.Track.Artists[0].Name + " - " + track.Track.TimeDuration().Round(time.Second).String(),
		})
	}
	return items, nil
}

func ArtistsView(ctx *gctx.Context, client *spotify.Client) ([]list.Item, error) {
	items := []list.Item{}
	artists, err := commands.UserArtists(ctx, client, 1)
	if err != nil {
		return nil, err
	}
	for _, artist := range artists.Artists {
		items = append(items, mainItem{
			Name:        artist.Name,
			ID:          artist.ID,
			Desc:        fmt.Sprintf("%d followers, genres: %s, popularity: %d", artist.Followers.Count, artist.Genres, artist.Popularity),
			SpotifyItem: artist.SimpleArtist,
		})
	}
	return items, nil
}

func SearchArtistsView(ctx *gctx.Context, client *spotify.Client, artists *spotify.FullArtistPage) ([]list.Item, error) {
	items := []list.Item{}
	for _, artist := range artists.Artists {
		items = append(items, mainItem{
			Name:        artist.Name,
			ID:          artist.ID,
			Desc:        fmt.Sprintf("%d followers, genres: %s, popularity: %d", artist.Followers.Count, artist.Genres, artist.Popularity),
			SpotifyItem: artist.SimpleArtist,
		})
	}
	return items, nil
}

func SearchView(ctx *gctx.Context, client *spotify.Client, search string) ([]list.Item, error) {
	items := []list.Item{}

	result, err := commands.Search(ctx, client, search, 1)
	if err != nil {
		return nil, err
	}
	items = append(items, mainItem{
		Name:        "Tracks",
		Desc:        "Search results",
		SpotifyItem: result.Tracks,
	})
	items = append(items, mainItem{
		Name:        "Albums",
		Desc:        "Search results",
		SpotifyItem: result.Albums,
	})
	items = append(items, mainItem{
		Name:        "Artists",
		Desc:        "Search results",
		SpotifyItem: result.Artists,
	})
	items = append(items, mainItem{
		Name:        "Playlists",
		Desc:        "Search results",
		SpotifyItem: result.Playlists,
	})
	return items, nil
}

func AlbumsView(ctx *gctx.Context, client *spotify.Client) ([]list.Item, error) {
	items := []list.Item{}
	albums, err := commands.UserAlbums(ctx, client, 1)
	if err != nil {
		return nil, err
	}
	for _, album := range albums.Albums {
		items = append(items, mainItem{
			Name:        album.Name,
			ID:          album.ID,
			Desc:        fmt.Sprintf("%s, %d tracks", album.Artists[0].Name, album.Tracks.Total),
			SpotifyItem: album.SimpleAlbum,
		})
	}
	return items, nil
}

func SearchPlaylistsView(ctx *gctx.Context, client *spotify.Client, playlists *spotify.SimplePlaylistPage) ([]list.Item, error) {
	items := []list.Item{}
	for _, playlist := range playlists.Playlists {
		items = append(items, mainItem{
			Name:        playlist.Name,
			Desc:        playlist.Description,
			SpotifyItem: playlist,
		})
	}
	return items, nil
}

func SearchAlbumsView(ctx *gctx.Context, client *spotify.Client, albums *spotify.SimpleAlbumPage) ([]list.Item, error) {
	items := []list.Item{}
	for _, album := range albums.Albums {
		items = append(items, mainItem{
			Name:        album.Name,
			ID:          album.ID,
			Desc:        fmt.Sprintf("%s, %d", album.Artists[0].Name, album.ReleaseDateTime()),
			SpotifyItem: album,
		})
	}
	return items, nil
}

func ArtistAlbumsView(ctx *gctx.Context, album spotify.ID, client *spotify.Client) ([]list.Item, error) {
	items := []list.Item{}
	albums, err := commands.ArtistAlbums(ctx, client, album, 1)
	if err != nil {
		return nil, err
	}
	for _, album := range albums.Albums {
		items = append(items, mainItem{
			Name:        album.Name,
			ID:          album.ID,
			Desc:        fmt.Sprintf("%s by %s", album.AlbumType, album.Artists[0].Name),
			SpotifyItem: album,
		})
	}
	return items, err
}

func AlbumTracksView(ctx *gctx.Context, album spotify.ID, client *spotify.Client) ([]list.Item, error) {
	items := []list.Item{}
	tracks, err := commands.AlbumTracks(ctx, client, album, 1)
	if err != nil {
		return nil, err
	}
	for _, track := range tracks.Tracks {
		items = append(items, mainItem{
			Name:     track.Name,
			Artist:   track.Artists[0],
			Duration: track.TimeDuration().Round(time.Second).String(),
			ID:       track.ID,
			Desc:     track.Artists[0].Name + " - " + track.TimeDuration().Round(time.Second).String(),
		})
	}
	return items, err
}

func SearchTracksView(ctx *gctx.Context, client *spotify.Client, tracks *spotify.FullTrackPage) ([]list.Item, error) {
	items := []list.Item{}
	for _, track := range tracks.Tracks {
		items = append(items, mainItem{
			Name:     track.Name,
			Artist:   track.Artists[0],
			Duration: track.TimeDuration().Round(time.Second).String(),
			ID:       track.ID,
			Desc:     track.Artists[0].Name + " - " + track.TimeDuration().Round(time.Second).String(),
		})
	}
	return items, nil
}

func SavedTracksView(ctx *gctx.Context, client *spotify.Client) ([]list.Item, error) {
	items := []list.Item{}
	tracks, err := commands.TrackList(ctx, client, 1)
	if err != nil {
		return nil, err
	}
	for _, track := range tracks.Tracks {
		items = append(items, mainItem{
			Name:     track.Name,
			Artist:   track.Artists[0],
			Duration: track.TimeDuration().Round(time.Second).String(),
			ID:       track.ID,
			Desc:     track.Artists[0].Name + " - " + track.TimeDuration().Round(time.Second).String(),
		})
	}
	return items, err
}

func MainView(ctx *gctx.Context, client *spotify.Client) ([]list.Item, error) {
	var wg sync.WaitGroup
	var saved_items *spotify.SavedTrackPage
	var playlists *spotify.SimplePlaylistPage
	var artists *spotify.FullArtistCursorPage
	var albums *spotify.SavedAlbumPage

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		saved_items, err = commands.TrackList(ctx, client, 1)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		playlists, err = commands.Playlists(ctx, client, 1)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		artists, err = commands.UserArtists(ctx, client, 1)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		albums, err = commands.UserAlbums(ctx, client, 1)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	wg.Wait()

	items := []list.Item{}
	items = append(items, mainItem{
		Name:        "Saved Tracks",
		Desc:        fmt.Sprintf("%d saved songs", saved_items.Total),
		SpotifyItem: saved_items,
	})
	items = append(items, mainItem{
		Name:        "Albums",
		Desc:        fmt.Sprintf("%d albums", albums.Total),
		SpotifyItem: albums,
	})
	items = append(items, mainItem{
		Name:        "Artists",
		Desc:        fmt.Sprintf("%d artists", artists.Total),
		SpotifyItem: artists,
	})
	for _, playlist := range playlists.Playlists {
		items = append(items, mainItem{
			Name:        playlist.Name,
			Desc:        playlist.Description,
			SpotifyItem: playlist,
		})
	}
	return items, nil
}

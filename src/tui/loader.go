package tui

import (
	"fmt"
	"time"

	"gitea.asdf.cafe/abs3nt/gospt/src/commands"

	"github.com/charmbracelet/bubbles/list"
)

func (m *mainModel) LoadMoreItems() {
	switch m.mode {
	case "artist":
		albums, err := commands.ArtistAlbums(m.ctx, m.client, m.artist.ID, (page + 1))
		page++
		if err != nil {
			return
		}
		items := []list.Item{}
		for _, album := range albums.Albums {
			items = append(items, mainItem{
				Name:        album.Name,
				ID:          album.ID,
				Desc:        fmt.Sprintf("%s by %s", album.AlbumType, album.Artists[0].Name),
				SpotifyItem: album,
			})
		}
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
		main_updates <- m
		return
	case "artists":
		artists, err := commands.UserArtists(m.ctx, m.client, (page + 1))
		page++
		if err != nil {
			return
		}
		items := []list.Item{}
		for _, artist := range artists.Artists {
			items = append(items, mainItem{
				Name:        artist.Name,
				ID:          artist.ID,
				Desc:        fmt.Sprintf("%d followers, genres: %s, popularity: %d", artist.Followers.Count, artist.Genres, artist.Popularity),
				SpotifyItem: artist.SimpleArtist,
			})
		}
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
		main_updates <- m
		return
	case "album":
		tracks, err := commands.AlbumTracks(m.ctx, m.client, m.album.ID, (page + 1))
		page++
		if err != nil {
			return
		}
		items := []mainItem{}
		for _, track := range tracks.Tracks {
			items = append(items, mainItem{
				Name:     track.Name,
				Artist:   track.Artists[0],
				Duration: track.TimeDuration().Round(time.Second).String(),
				ID:       track.ID,
				Desc:     track.Artists[0].Name + " - " + track.TimeDuration().Round(time.Second).String(),
			})
		}
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
		main_updates <- m
		return
	case "albums":
		albums, err := commands.UserAlbums(m.ctx, m.client, (page + 1))
		page++
		if err != nil {
			return
		}
		items := []list.Item{}
		for _, album := range albums.Albums {
			items = append(items, mainItem{
				Name:        album.Name,
				ID:          album.ID,
				Desc:        fmt.Sprintf("%s, %d tracks", album.Artists[0].Name, album.Tracks.Total),
				SpotifyItem: album.SimpleAlbum,
			})
		}
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
		main_updates <- m
		return
	case "main":
		playlists, err := commands.Playlists(m.ctx, m.client, (page + 1))
		page++
		if err != nil {
			return
		}
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
		main_updates <- m
		return
	case "playlist":
		tracks, err := commands.PlaylistTracks(m.ctx, m.client, m.playlist.ID, (page + 1))
		page++
		if err != nil {
			return
		}
		items := []mainItem{}
		for _, track := range tracks.Tracks {
			items = append(items, mainItem{
				Name:     track.Track.Name,
				Artist:   track.Track.Artists[0],
				Duration: track.Track.TimeDuration().Round(time.Second).String(),
				ID:       track.Track.ID,
				Desc:     track.Track.Artists[0].Name + " - " + track.Track.TimeDuration().Round(time.Second).String(),
			})
		}
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
		main_updates <- m
		return
	case "tracks":
		tracks, err := commands.TrackList(m.ctx, m.client, (page + 1))
		page++
		if err != nil {
			return
		}
		page++
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
		for _, item := range items {
			m.list.InsertItem(len(m.list.Items())+1, item)
		}
		main_updates <- m
		return
	}
}

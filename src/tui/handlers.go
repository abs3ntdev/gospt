package tui

import (
	"fmt"

	"gospt/src/commands"
	"gospt/src/gctx"

	"github.com/zmb3/spotify/v2"
)

func HandlePlayWithContext(ctx *gctx.Context, client *spotify.Client, uri *spotify.URI, pos int) {
	var err error
	err = commands.PlaySongInPlaylist(ctx, client, uri, pos)
	if err != nil {
		return
	}
}

func HandleRadio(ctx *gctx.Context, client *spotify.Client, id spotify.ID) {
	err := commands.RadioGivenSong(ctx, client, id, 0)
	if err != nil {
		return
	}
}

func HandleAlbumRadio(ctx *gctx.Context, client *spotify.Client, id spotify.ID) {
	err := commands.RadioFromAlbum(ctx, client, id)
	if err != nil {
		return
	}
}

func HandleSeek(ctx *gctx.Context, client *spotify.Client, fwd bool) {
	err := commands.Seek(ctx, client, fwd)
	if err != nil {
		return
	}
}

func HandleVolume(ctx *gctx.Context, client *spotify.Client, up bool) {
	vol := 10
	if !up {
		vol = -10
	}
	err := commands.ChangeVolume(ctx, client, vol)
	if err != nil {
		return
	}
}

func HandleArtistRadio(ctx *gctx.Context, client *spotify.Client, id spotify.ID) {
	err := commands.RadioGivenArtist(ctx, client, id)
	if err != nil {
		return
	}
}

func HandleAlbumArtist(ctx *gctx.Context, client *spotify.Client, id spotify.ID) {
	err := commands.RadioGivenArtist(ctx, client, id)
	if err != nil {
		return
	}
}

func HandlePlaylistRadio(ctx *gctx.Context, client *spotify.Client, playlist spotify.SimplePlaylist) {
	err := commands.RadioFromPlaylist(ctx, client, playlist)
	if err != nil {
		return
	}
}

func HandleLibraryRadio(ctx *gctx.Context, client *spotify.Client) {
	err := commands.RadioFromSavedTracks(ctx, client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func HandlePlayLikedSong(ctx *gctx.Context, client *spotify.Client, position int) {
	err := commands.PlayLikedSongs(ctx, client, position)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func HandlePlayTrack(ctx *gctx.Context, client *spotify.Client, track spotify.ID) {
	err := commands.QueueSong(ctx, client, track)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = commands.Next(ctx, client, 1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func HandleSetDevice(ctx *gctx.Context, client *spotify.Client, player spotify.PlayerDevice) {
	var err error
	err = commands.SetDevice(ctx, client, player)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

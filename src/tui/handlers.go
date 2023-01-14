package tui

import (
	"fmt"

	"gospt/src/commands"
	"gospt/src/gctx"

	"github.com/zmb3/spotify/v2"
)

func HandlePlay(ctx *gctx.Context, client *spotify.Client, uri *spotify.URI, pos int) {
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

func HandleAlbumArtist(ctx *gctx.Context, client *spotify.Client, id spotify.ID) {
	err := commands.RadioGivenArtist(ctx, client, id)
	if err != nil {
		return
	}
}

func HandlePlaylistRadio(ctx *gctx.Context, client *spotify.Client, playlist spotify.SimplePlaylist) {
	err := commands.RadioFromPlaylist(ctx, client, playlist)
	if err != nil {
		fmt.Println("AHHHHHHHHHHHHHHHHHH", err.Error())
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

func HandleSetDevice(ctx *gctx.Context, client *spotify.Client, player spotify.PlayerDevice) {
	fmt.Println("WHOA")
	var err error
	err = commands.SetDevice(ctx, client, player)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

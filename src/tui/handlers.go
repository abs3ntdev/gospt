package tui

import (
	"github.com/zmb3/spotify/v2"

	"git.asdf.cafe/abs3nt/gospt/src/commands"
	"git.asdf.cafe/abs3nt/gospt/src/gctx"
)

func HandlePlayWithContext(ctx *gctx.Context, commands *commands.Commands, uri *spotify.URI, pos *int) {
	err := commands.PlaySongInPlaylist(ctx, uri, pos)
	if err != nil {
		return
	}
}

func HandleRadio(ctx *gctx.Context, commands *commands.Commands, song spotify.SimpleTrack) {
	err := commands.RadioGivenSong(ctx, song, 0)
	if err != nil {
		return
	}
}

func HandleAlbumRadio(ctx *gctx.Context, commands *commands.Commands, album spotify.SimpleAlbum) {
	err := commands.RadioFromAlbum(ctx, album)
	if err != nil {
		return
	}
}

func HandleSeek(ctx *gctx.Context, commands *commands.Commands, fwd bool) {
	err := commands.Seek(ctx, fwd)
	if err != nil {
		return
	}
}

func HandleVolume(ctx *gctx.Context, commands *commands.Commands, up bool) {
	vol := 10
	if !up {
		vol = -10
	}
	err := commands.ChangeVolume(ctx, vol)
	if err != nil {
		return
	}
}

func HandleArtistRadio(ctx *gctx.Context, commands *commands.Commands, artist spotify.SimpleArtist) {
	err := commands.RadioGivenArtist(ctx, artist)
	if err != nil {
		return
	}
}

func HandleAlbumArtist(ctx *gctx.Context, commands *commands.Commands, artist spotify.SimpleArtist) {
	err := commands.RadioGivenArtist(ctx, artist)
	if err != nil {
		return
	}
}

func HandlePlaylistRadio(ctx *gctx.Context, commands *commands.Commands, playlist spotify.SimplePlaylist) {
	err := commands.RadioFromPlaylist(ctx, playlist)
	if err != nil {
		return
	}
}

func HandleLibraryRadio(ctx *gctx.Context, commands *commands.Commands) {
	err := commands.RadioFromSavedTracks(ctx)
	if err != nil {
		return
	}
}

func HandlePlayLikedSong(ctx *gctx.Context, commands *commands.Commands, position int) {
	err := commands.PlayLikedSongs(ctx, position)
	if err != nil {
		return
	}
}

func HandlePlayTrack(ctx *gctx.Context, commands *commands.Commands, track spotify.ID) {
	err := commands.QueueSong(ctx, track)
	if err != nil {
		return
	}
	err = commands.Next(ctx, 1, false)
	if err != nil {
		return
	}
}

func HandleNextInQueue(ctx *gctx.Context, commands *commands.Commands, amt int) {
	err := commands.Next(ctx, amt, true)
	if err != nil {
		return
	}
}

func HandleQueueItem(ctx *gctx.Context, commands *commands.Commands, item spotify.ID) {
	err := commands.QueueSong(ctx, item)
	if err != nil {
		return
	}
}

func HandleDeleteTrackFromPlaylist(ctx *gctx.Context, commands *commands.Commands, item, playlist spotify.ID) {
	err := commands.DeleteTracksFromPlaylist(ctx, []spotify.ID{item}, playlist)
	if err != nil {
		return
	}
}

func HandleSetDevice(ctx *gctx.Context, commands *commands.Commands, player spotify.PlayerDevice) {
	err := commands.SetDevice(ctx, player)
	if err != nil {
		return
	}
}

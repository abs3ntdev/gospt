package commands

import (
	"encoding/json"
	"fmt"

	"gospt/ctx"

	"github.com/zmb3/spotify/v2"
)

func Play(ctx *ctx.Context, client *spotify.Client) error {
	var err error
	err = client.Play(ctx)
	if err != nil {
		return err
	}
	ctx.Println("Playing!")
	return nil
}

func Pause(ctx *ctx.Context, client *spotify.Client) error {
	err := client.Pause(ctx)
	if err != nil {
		return err
	}
	ctx.Println("Pausing!")
	return nil
}

func Skip(ctx *ctx.Context, client *spotify.Client) error {
	err := client.Next(ctx)
	if err != nil {
		return err
	}
	ctx.Println("Skipping!")
	return nil
}

func Status(ctx *ctx.Context, client *spotify.Client) error {
	state, err := client.PlayerState(ctx)
	if err != nil {
		return err
	}
	return PrintState(state)
}

func Shuffle(ctx *ctx.Context, client *spotify.Client) error {
	state, err := client.PlayerState(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get current playstate")
	}
	err = client.Shuffle(ctx, !state.ShuffleState)
	if err != nil {
		return err
	}
	ctx.Println("Shuffle set to", !state.ShuffleState)
	return nil
}

func TrackList(ctx *ctx.Context, client *spotify.Client, page int) (*spotify.SavedTrackPage, error) {
	return client.CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func PrintState(state *spotify.PlayerState) error {
	state.Item.AvailableMarkets = []string{}
	state.Item.Album.AvailableMarkets = []string{}
	out, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

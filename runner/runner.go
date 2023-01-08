package runner

import (
	"fmt"

	"gospt/ctx"
	"gospt/tui"

	"github.com/zmb3/spotify/v2"
)

func Run(ctx *ctx.Context, client *spotify.Client, args []string) error {
	if len(args) == 0 {
		user, err := client.CurrentUser(ctx)
		if err != nil {
			return fmt.Errorf("Failed to get current user")
		}
		fmt.Println("The following commands are currently supported:\nplay pause next shuffle\nhave fun", user.DisplayName)
		return nil
	}
	switch args[0] {
	case "play":
		return Play(ctx, client)
	case "pause":
		return Pause(ctx, client)
	case "next":
		return Skip(ctx, client)
	case "shuffle":
		return Shuffle(ctx, client)
	case "tracks":
		return GetTracks(ctx, client, args)
	default:
		return fmt.Errorf("Unsupported Command")
	}
}

func Play(ctx *ctx.Context, client *spotify.Client) error {
	err := client.Play(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Playing!")
	return nil
}

func Pause(ctx *ctx.Context, client *spotify.Client) error {
	err := client.Pause(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Pausing!")
	return nil
}

func Skip(ctx *ctx.Context, client *spotify.Client) error {
	err := client.Next(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Skipping!")
	return nil
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
	fmt.Println("Shuffle set to", !state.ShuffleState)
	return nil
}

func GetTracks(ctx *ctx.Context, client *spotify.Client, args []string) error {
	return tui.DisplayList(ctx, client)
}

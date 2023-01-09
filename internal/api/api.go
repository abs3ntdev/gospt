package api

import (
	"fmt"

	"gospt/internal/commands"
	"gospt/internal/gctx"
	"gospt/internal/tui"

	"github.com/zmb3/spotify/v2"
)

func Run(ctx *gctx.Context, client *spotify.Client, args []string) error {
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
		return commands.Play(ctx, client)
	case "pause":
		return commands.Pause(ctx, client)
	case "toggleplay":
		return commands.TogglePlay(ctx, client)
	case "next":
		return commands.Skip(ctx, client)
	case "previous":
		return commands.Previous(ctx, client)
	case "playurl":
		return commands.PlayUrl(ctx, client, args)
	case "like":
		return commands.Like(ctx, client)
	case "unlike":
		return commands.Unlike(ctx, client)
	case "shuffle":
		return commands.Shuffle(ctx, client)
	case "repeat":
		return commands.Repeat(ctx, client)
	case "radio":
		return commands.Radio(ctx, client)
	case "clearradio":
		return commands.ClearRadio(ctx, client)
	case "refillradio":
		return commands.RefillRadio(ctx, client)
	case "tracks":
		return tui.DisplayList(ctx, client)
	case "playlists":
		return tui.DisplayPlaylists(ctx, client)
	case "status":
		return commands.Status(ctx, client)
	case "devices":
		return commands.Devices(ctx, client)
	case "setdevice":
		return tui.DisplayDevices(ctx, client)
	default:
		return fmt.Errorf("Unsupported Command")
	}
}

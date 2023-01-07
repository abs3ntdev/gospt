package runner

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func Run(client *spotify.Client, args []string) error {
	if len(args) == 0 {
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			return fmt.Errorf("Failed to get current user")
		}
		fmt.Println("YOU ARE", user.DisplayName)
		return nil
	}
	ctx := context.Background()
	switch args[0] {
	case "play":
		err := client.Play(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Playing!")
	case "pause":
		err := client.Pause(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Pausing!")
	case "next":
		err := client.Next(ctx)
		if err != nil {
			return err
		}
		fmt.Println("Skipping!")
	case "shuffle":
		state, err := client.PlayerState(ctx)
		if err != nil {
			return fmt.Errorf("Failed to get current playstate")
		}
		err = client.Shuffle(ctx, !state.ShuffleState)
		if err != nil {
			return err
		}
		fmt.Println("Shuffle set to", !state.ShuffleState)
	}
	return nil
}

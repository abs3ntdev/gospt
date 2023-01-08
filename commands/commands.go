package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gospt/ctx"

	"github.com/zmb3/spotify/v2"
)

func Play(ctx *ctx.Context, client *spotify.Client) error {
	err := client.Play(ctx)
	if err != nil {
		if isNoActiveError(err) {
			return playWithTransfer(ctx, client)
		}
		return err
	}
	ctx.Println("Playing!")
	return nil
}

func PlayUrl(ctx *ctx.Context, client *spotify.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Please provide a url")
	}
	url, err := url.Parse(args[1])
	if err != nil {
		return err
	}
	track_id := strings.Split(url.Path, "/")[2]
	err = client.QueueSong(ctx, spotify.ID(track_id))
	if err != nil {
		if isNoActiveError(err) {
			return queueWithTransfer(ctx, client, spotify.ID(track_id))
		}
		return err
	}
	err = client.Next(ctx)
	ctx.Println("Playing!")
	return nil
}

func Devices(ctx *ctx.Context, client *spotify.Client) error {
	devices, err := client.PlayerDevices(ctx)
	if err != nil {
		return err
	}
	return PrintDevices(devices)
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

func PrintDevices(devices []spotify.PlayerDevice) error {
	out, err := json.MarshalIndent(devices, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func SetDevice(ctx *ctx.Context, client *spotify.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Please provide your device ID")
	}
	devices, err := client.PlayerDevices(ctx)
	if err != nil {
		return err
	}
	var set_device spotify.PlayerDevice
	for _, device := range devices {
		if device.ID.String() == args[1] {
			set_device = device
			break
		}
	}
	out, err := json.MarshalIndent(set_device, "", " ")
	if err != nil {
		return err
	}
	homdir, _ := os.UserHomeDir()
	err = ioutil.WriteFile(filepath.Join(homdir, ".config/gospt/device.json"), out, 0o644)
	if err != nil {
		return err
	}
	fmt.Println("Your device has been set to: ", set_device.Name)
	return nil
}

func isNoActiveError(err error) bool {
	return strings.Contains(err.Error(), "No active device found")
}

func playWithTransfer(ctx *ctx.Context, client *spotify.Client) error {
	configDir, _ := os.UserConfigDir()
	deviceFile, err := os.Open(filepath.Join(configDir, "gospt/device.json"))
	if err != nil {
		return err
	}
	defer deviceFile.Close()
	deviceValue, err := ioutil.ReadAll(deviceFile)
	if err != nil {
		return err
	}
	var device *spotify.PlayerDevice
	err = json.Unmarshal(deviceValue, &device)
	if err != nil {
		return err
	}
	err = client.TransferPlayback(ctx, device.ID, true)
	if err != nil {
		return err
	}
	ctx.Println("Playing!")
	return nil
}

func queueWithTransfer(ctx *ctx.Context, client *spotify.Client, track_id spotify.ID) error {
	configDir, _ := os.UserConfigDir()
	deviceFile, err := os.Open(filepath.Join(configDir, "gospt/device.json"))
	if err != nil {
		return err
	}
	defer deviceFile.Close()
	deviceValue, err := ioutil.ReadAll(deviceFile)
	if err != nil {
		return err
	}
	var device *spotify.PlayerDevice
	err = json.Unmarshal(deviceValue, &device)
	if err != nil {
		return err
	}
	err = client.TransferPlayback(ctx, device.ID, true)
	if err != nil {
		return err
	}
	err = client.QueueSong(ctx, track_id)
	if err != nil {
		return err
	}
	err = client.Next(ctx)
	if err != nil {
		return err
	}
	ctx.Println("Playing!")
	return nil
}

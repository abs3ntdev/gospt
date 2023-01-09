package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gospt/internal/gctx"

	"github.com/zmb3/spotify/v2"
)

func Play(ctx *gctx.Context, client *spotify.Client) error {
	err := client.Play(ctx)
	if err != nil {
		if isNoActiveError(err) {
			err := activateDevice(ctx, client)
			if err != nil {
				return err
			}
			err = client.Play(ctx)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func PlayUrl(ctx *gctx.Context, client *spotify.Client, args []string) error {
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
			err := activateDevice(ctx, client)
			if err != nil {
				return err
			}
			err = client.QueueSong(ctx, spotify.ID(track_id))
			if err != nil {
				return err
			}
			err = client.Next(ctx)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	err = client.Next(ctx)
	if err != nil {
		return err
	}
	return nil
}

func QueueSong(ctx *gctx.Context, client *spotify.Client, id spotify.ID) error {
	err := client.QueueSong(ctx, id)
	if err != nil {
		if isNoActiveError(err) {
			err := activateDevice(ctx, client)
			if err != nil {
				return err
			}
			err = client.QueueSong(ctx, id)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	return nil
}

func RadioGivenSong(ctx *gctx.Context, client *spotify.Client, song_id spotify.ID) error {
	seed := spotify.Seeds{
		Tracks: []spotify.ID{song_id},
	}
	recomendations, err := client.GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(99))
	if err != nil {
		return err
	}
	recomendationIds := []spotify.ID{}
	for _, song := range recomendations.Tracks {
		recomendationIds = append(recomendationIds, song.ID)
	}
	err = ClearRadio(ctx, client)
	if err != nil {
		return err
	}
	radioPlaylist, err := GetRadioPlaylist(ctx, client)
	if err != nil {
		return err
	}
	queue := []spotify.ID{song_id}
	queue = append(queue, recomendationIds...)
	_, err = client.AddTracksToPlaylist(ctx, radioPlaylist.ID, queue...)
	if err != nil {
		return err
	}
	client.PlayOpt(ctx, &spotify.PlayOptions{
		PlaybackContext: &radioPlaylist.URI,
		PlaybackOffset: &spotify.PlaybackOffset{
			Position: 0,
		},
	})
	err = client.Repeat(ctx, "context")
	if err != nil {
		return err
	}
	for i := 0; i < 4; i++ {
		id := rand.Intn(len(recomendationIds)-2) + 1
		seed := spotify.Seeds{
			Tracks: []spotify.ID{recomendationIds[id]},
		}
		additional_recs, err := client.GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
		if err != nil {
			return err
		}
		additionalRecsIds := []spotify.ID{}
		for _, song := range additional_recs.Tracks {
			additionalRecsIds = append(additionalRecsIds, song.ID)
		}
		_, err = client.AddTracksToPlaylist(ctx, radioPlaylist.ID, additionalRecsIds...)
		if err != nil {
			return err
		}
	}
	return nil
}

func Radio(ctx *gctx.Context, client *spotify.Client) error {
	rand.Seed(time.Now().Unix())
	current_song, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	var seed_song spotify.SimpleTrack

	if current_song.Item != nil {
		seed_song = current_song.Item.SimpleTrack
	}
	if current_song.Item == nil {
		err := activateDevice(ctx, client)
		if err != nil {
			return err
		}
		tracks, err := client.CurrentUsersTracks(ctx, spotify.Limit(10))
		if err != nil {
			return err
		}
		seed_song = tracks.Tracks[rand.Intn(len(tracks.Tracks))].SimpleTrack
	} else {
		if !current_song.Playing {
			tracks, err := client.CurrentUsersTracks(ctx, spotify.Limit(10))
			if err != nil {
				return err
			}
			seed_song = tracks.Tracks[rand.Intn(len(tracks.Tracks))].SimpleTrack
		}
	}
	return RadioGivenSong(ctx, client, seed_song.ID)
}

func RefillRadio(ctx *gctx.Context, client *spotify.Client) error {
	status, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	to_remove := []spotify.ID{}
	radioPlaylist, err := GetRadioPlaylist(ctx, client)
	found := false
	page := 0
	for !found {
		tracks, err := client.GetPlaylistItems(ctx, radioPlaylist.ID, spotify.Limit(50), spotify.Offset(page*50))
		if err != nil {
			return err
		}
		for _, track := range tracks.Items {
			if track.Track.Track.ID == status.Item.ID {
				found = true
				break
			}
			to_remove = append(to_remove, track.Track.Track.ID)
		}
		page++
	}
	recomendationIds := []spotify.ID{}
	if len(to_remove) > 0 {
		_, err = client.RemoveTracksFromPlaylist(ctx, radioPlaylist.ID, to_remove...)
		if err != nil {
			return err
		}
		current_song, err := client.PlayerCurrentlyPlaying(ctx)
		if err != nil {
			return err
		}
		seed := spotify.Seeds{
			Tracks: []spotify.ID{current_song.Item.ID},
		}
		recomendations, err := client.GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
		if err != nil {
			return err
		}
		for idx, song := range recomendations.Tracks {
			if idx >= len(to_remove) {
				break
			}
			recomendationIds = append(recomendationIds, song.ID)
		}
		_, err = client.AddTracksToPlaylist(ctx, radioPlaylist.ID, recomendationIds...)
		if err != nil {
			return err
		}
	}
	return nil
}

func ClearRadio(ctx *gctx.Context, client *spotify.Client) error {
	radioPlaylist, err := GetRadioPlaylist(ctx, client)
	if err != nil {
		return err
	}
	err = client.UnfollowPlaylist(ctx, radioPlaylist.ID)
	if err != nil {
		return err
	}
	configDir, _ := os.UserConfigDir()
	os.Remove(filepath.Join(configDir, "gospt/radio.json"))
	client.Pause(ctx)
	return nil
}

func Devices(ctx *gctx.Context, client *spotify.Client) error {
	devices, err := client.PlayerDevices(ctx)
	if err != nil {
		return err
	}
	return PrintDevices(devices)
}

func Pause(ctx *gctx.Context, client *spotify.Client) error {
	err := client.Pause(ctx)
	if err != nil {
		return err
	}
	return nil
}

func TogglePlay(ctx *gctx.Context, client *spotify.Client) error {
	current, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	if !current.Playing {
		return Play(ctx, client)
	}
	return Pause(ctx, client)
}

func Like(ctx *gctx.Context, client *spotify.Client) error {
	playing, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	err = client.AddTracksToLibrary(ctx, playing.Item.ID)
	if err != nil {
		return err
	}
	return nil
}

func Unlike(ctx *gctx.Context, client *spotify.Client) error {
	playing, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	err = client.RemoveTracksFromLibrary(ctx, playing.Item.ID)
	if err != nil {
		return err
	}
	return nil
}

func Skip(ctx *gctx.Context, client *spotify.Client) error {
	err := client.Next(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Previous(ctx *gctx.Context, client *spotify.Client) error {
	err := client.Previous(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Status(ctx *gctx.Context, client *spotify.Client) error {
	state, err := client.PlayerState(ctx)
	if err != nil {
		return err
	}
	return PrintState(state)
}

func Shuffle(ctx *gctx.Context, client *spotify.Client) error {
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

func Repeat(ctx *gctx.Context, client *spotify.Client) error {
	state, err := client.PlayerState(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get current playstate")
	}
	fmt.Println(state.RepeatState)
	newState := "off"
	if state.RepeatState == "off" {
		newState = "context"
	}
	// spotifyd only supports binary value for repeat, context or off, change when/if spotifyd is better
	err = client.Repeat(ctx, newState)
	if err != nil {
		return err
	}
	ctx.Println("Repeat set to", newState)
	return nil
}

func TrackList(ctx *gctx.Context, client *spotify.Client, page int) (*spotify.SavedTrackPage, error) {
	return client.CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func Playlists(ctx *gctx.Context, client *spotify.Client, page int) (*spotify.SimplePlaylistPage, error) {
	return client.CurrentUsersPlaylists(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func PlaylistTracks(ctx *gctx.Context, client *spotify.Client, playlist spotify.ID, page int) (*spotify.PlaylistTrackPage, error) {
	return client.GetPlaylistTracks(ctx, playlist, spotify.Limit(50), spotify.Offset((page-1)*50))
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

func SetDevice(ctx *gctx.Context, client *spotify.Client, device spotify.PlayerDevice) error {
	out, err := json.MarshalIndent(device, "", " ")
	if err != nil {
		return err
	}
	configDir, _ := os.UserConfigDir()
	err = ioutil.WriteFile(filepath.Join(configDir, "gospt/device.json"), out, 0o644)
	if err != nil {
		return err
	}
	err = activateDevice(ctx, client)
	if err != nil {
		return err
	}
	fmt.Println("Your device has been set to: ", device.Name)
	return nil
}

func isNoActiveError(err error) bool {
	return strings.Contains(err.Error(), "No active device found")
}

func activateDevice(ctx *gctx.Context, client *spotify.Client) error {
	to_play := true
	current, err := client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	if current.Item == nil || !current.Playing {
		to_play = false
	}
	configDir, _ := os.UserConfigDir()
	if _, err := os.Stat(filepath.Join(configDir, "gospt/device.json")); err == nil {
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
		err = client.TransferPlayback(ctx, device.ID, to_play)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("YOU MUST RUN gospt setdevice FIRST")
	}
	return nil
}

func GetRadioPlaylist(ctx *gctx.Context, client *spotify.Client) (*spotify.FullPlaylist, error) {
	configDir, _ := os.UserConfigDir()
	if _, err := os.Stat(filepath.Join(configDir, "gospt/radio.json")); err == nil {
		playlistFile, err := os.Open(filepath.Join(configDir, "gospt/radio.json"))
		if err != nil {
			return nil, err
		}
		defer playlistFile.Close()
		playlistValue, err := ioutil.ReadAll(playlistFile)
		if err != nil {
			return nil, err
		}
		var playlist *spotify.FullPlaylist
		err = json.Unmarshal(playlistValue, &playlist)
		if err != nil {
			return nil, err
		}
		return playlist, nil
	}
	// private flag doesnt work
	playlist, err := client.CreatePlaylistForUser(ctx, ctx.UserId, "Radio", "Automanaged radio playlist", false, false)
	if err != nil {
		return nil, err
	}
	out, err := json.MarshalIndent(playlist, "", " ")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filepath.Join(configDir, "gospt/radio.json"), out, 0o644)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

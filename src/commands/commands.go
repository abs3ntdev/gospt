package commands

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gfx.cafe/util/go/frand"
	"git.asdf.cafe/abs3nt/gospt/src/auth"
	"git.asdf.cafe/abs3nt/gospt/src/cache"
	"git.asdf.cafe/abs3nt/gospt/src/gctx"
	"git.asdf.cafe/abs3nt/gospt/src/youtube"

	"github.com/zmb3/spotify/v2"
	_ "modernc.org/sqlite"
)

type Commands struct {
	Context *gctx.Context
	cl      *spotify.Client
	mu      sync.RWMutex

	user string
}

func (c *Commands) Client() *spotify.Client {
	c.mu.Lock()
	if c.cl == nil {
		c.cl = c.connectClient()
	}
	c.mu.Unlock()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cl
}

func (c *Commands) User() string {
	c.Client()
	return c.user
}

func (c *Commands) connectClient() *spotify.Client {
	ctx := c.Context
	client, err := auth.GetClient(ctx)
	if err != nil {
		panic(err)
	}
	currentUser, err := client.CurrentUser(ctx)
	if err != nil {
		panic(err)
	}
	c.user = currentUser.ID
	return client
}

func (c *Commands) SetVolume(ctx *gctx.Context, vol int) error {
	return c.Client().Volume(ctx, vol)
}

func (c *Commands) SetPosition(ctx *gctx.Context, pos int) error {
	err := c.Client().Seek(ctx, pos)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Seek(ctx *gctx.Context, fwd bool) error {
	current, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	newPos := current.Progress + 5000
	if !fwd {
		newPos = current.Progress - 5000
	}
	err = c.Client().Seek(ctx, newPos)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) ChangeVolume(ctx *gctx.Context, vol int) error {
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return err
	}
	newVolume := state.Device.Volume + vol
	if newVolume > 100 {
		newVolume = 100
	}
	if newVolume < 0 {
		newVolume = 0
	}
	return c.Client().Volume(ctx, newVolume)
}

func (c *Commands) Play(ctx *gctx.Context) error {
	err := c.Client().Play(ctx)
	if err != nil {
		if isNoActiveError(err) {
			deviceID, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
				DeviceID: &deviceID,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (c *Commands) ActiveDeviceExists(ctx *gctx.Context) bool {
	current, err := c.Client().PlayerDevices(ctx)
	if err != nil {
		return false
	}
	for _, dev := range current {
		if dev.Active {
			return true
		}
	}
	return false
}

func (c *Commands) UserArtists(ctx *gctx.Context, page int) (*spotify.FullArtistCursorPage, error) {
	artists, err := c.Client().CurrentUsersFollowedArtists(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
	if err != nil {
		return nil, err
	}
	return artists, nil
}

func (c *Commands) ArtistAlbums(ctx *gctx.Context, artist spotify.ID, page int) (*spotify.SimpleAlbumPage, error) {
	albums, err := c.Client().GetArtistAlbums(ctx, artist, []spotify.AlbumType{1, 2, 3, 4}, spotify.Market(spotify.CountryUSA), spotify.Limit(50), spotify.Offset((page-1)*50))
	if err != nil {
		return nil, err
	}
	return albums, nil
}

func (c *Commands) Search(ctx *gctx.Context, search string, page int) (*spotify.SearchResult, error) {
	result, err := c.Client().Search(ctx, search, spotify.SearchTypeAlbum|spotify.SearchTypeArtist|spotify.SearchTypeTrack|spotify.SearchTypePlaylist, spotify.Limit(50), spotify.Offset((page-1)*50))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Commands) AlbumTracks(ctx *gctx.Context, album spotify.ID, page int) (*spotify.SimpleTrackPage, error) {
	tracks, err := c.Client().GetAlbumTracks(ctx, album, spotify.Limit(50), spotify.Offset((page-1)*50), spotify.Market(spotify.CountryUSA))
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (c *Commands) UserAlbums(ctx *gctx.Context, page int) (*spotify.SavedAlbumPage, error) {
	return c.Client().CurrentUsersAlbums(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func (c *Commands) UserQueue(ctx *gctx.Context) (*spotify.Queue, error) {
	return c.Client().GetQueue(ctx)
}

func (c *Commands) PlayUrl(ctx *gctx.Context, args []string) error {
	url, err := url.Parse(args[0])
	if err != nil {
		return err
	}
	track_id := strings.Split(url.Path, "/")[2]
	err = c.Client().QueueSong(ctx, spotify.ID(track_id))
	if err != nil {
		if isNoActiveError(err) {
			deviceID, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().QueueSongOpt(ctx, spotify.ID(track_id), &spotify.PlayOptions{
				DeviceID: &deviceID,
			})
			if err != nil {
				return err
			}
			err = c.Client().NextOpt(ctx, &spotify.PlayOptions{
				DeviceID: &deviceID,
			})
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	err = c.Client().Next(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) QueueSong(ctx *gctx.Context, id spotify.ID) error {
	err := c.Client().QueueSong(ctx, id)
	if err != nil {
		if isNoActiveError(err) {
			deviceID, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().QueueSongOpt(ctx, id, &spotify.PlayOptions{
				DeviceID: &deviceID,
			})
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

func (c *Commands) PlaySongInPlaylist(ctx *gctx.Context, context *spotify.URI, offset int) error {
	e := c.Client().PlayOpt(ctx, &spotify.PlayOptions{
		PlaybackOffset:  &spotify.PlaybackOffset{Position: offset},
		PlaybackContext: context,
	})
	if e != nil {
		if isNoActiveError(e) {
			deviceID, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
				PlaybackOffset:  &spotify.PlaybackOffset{Position: offset},
				PlaybackContext: context,
				DeviceID:        &deviceID,
			})

			if err != nil {
				if isNoActiveError(err) {
					deviceID, err := c.activateDevice(ctx)
					if err != nil {
						return err
					}
					err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
						PlaybackOffset:  &spotify.PlaybackOffset{Position: offset},
						PlaybackContext: context,
						DeviceID:        &deviceID,
					})
					if err != nil {
						return err
					}
				}
			}
			err = c.Client().Play(ctx)
			if err != nil {
				return err
			}
		} else {
			return e
		}
	}
	return nil
}

func (c *Commands) PlayLikedSongs(ctx *gctx.Context, position int) error {
	err := c.ClearRadio(ctx)
	if err != nil {
		return err
	}
	playlist, _, err := c.GetRadioPlaylist(ctx, "Saved Songs")
	if err != nil {
		return err
	}
	songs, err := c.Client().CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset(position))
	if err != nil {
		return err
	}
	to_add := []spotify.ID{}
	for _, song := range songs.Tracks {
		to_add = append(to_add, song.ID)
	}
	_, err = c.Client().AddTracksToPlaylist(ctx, playlist.ID, to_add...)
	if err != nil {
		return err
	}
	err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
		PlaybackContext: &playlist.URI,
		PlaybackOffset: &spotify.PlaybackOffset{
			Position: 0,
		},
	})
	if err != nil {
		if isNoActiveError(err) {
			deviceID, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
				PlaybackContext: &playlist.URI,
				PlaybackOffset: &spotify.PlaybackOffset{
					Position: 0,
				},
				DeviceID: &deviceID,
			})
			if err != nil {
				return err
			}
		}
	}
	for page := 2; page <= 5; page++ {
		songs, err := c.Client().CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset((50*(page-1))+position))
		if err != nil {
			return err
		}
		to_add := []spotify.ID{}
		for _, song := range songs.Tracks {
			to_add = append(to_add, song.ID)
		}
		c.Client().AddTracksToPlaylist(ctx, playlist.ID, to_add...)
	}

	return err
}

func (c *Commands) RadioGivenArtist(ctx *gctx.Context, artist spotify.SimpleArtist) error {
	seed := spotify.Seeds{
		Artists: []spotify.ID{artist.ID},
	}
	recomendations, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
	if err != nil {
		return err
	}
	recomendationIds := []spotify.ID{}
	for _, song := range recomendations.Tracks {
		recomendationIds = append(recomendationIds, song.ID)
	}
	err = c.ClearRadio(ctx)
	if err != nil {
		return err
	}
	radioPlaylist, db, err := c.GetRadioPlaylist(ctx, artist.Name)
	if err != nil {
		return err
	}
	queue := []spotify.ID{}
	for _, rec := range recomendationIds {
		exists, err := c.SongExists(db, rec)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(rec)))
			if err != nil {
				return err
			}
			queue = append(queue, rec)
		}
	}
	_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, queue...)
	if err != nil {
		return err
	}
	c.Client().PlayOpt(ctx, &spotify.PlayOptions{
		PlaybackContext: &radioPlaylist.URI,
		PlaybackOffset: &spotify.PlaybackOffset{
			Position: 0,
		},
	})
	err = c.Client().Repeat(ctx, "context")
	if err != nil {
		return err
	}
	for i := 0; i < 4; i++ {
		id := frand.Intn(len(recomendationIds)-2) + 1
		seed := spotify.Seeds{
			Tracks: []spotify.ID{recomendationIds[id]},
		}
		additional_recs, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
		if err != nil {
			return err
		}
		additionalRecsIds := []spotify.ID{}
		for _, song := range additional_recs.Tracks {
			exists, err := c.SongExists(db, song.ID)
			if err != nil {
				return err
			}
			if !exists {
				_, err = db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(song.ID)))
				if err != nil {
					return err
				}
				additionalRecsIds = append(additionalRecsIds, song.ID)
			}
		}
		_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, additionalRecsIds...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) RadioGivenSong(ctx *gctx.Context, song spotify.SimpleTrack, pos int) error {
	start := time.Now().UnixMilli()
	seed := spotify.Seeds{
		Tracks: []spotify.ID{song.ID},
	}
	recomendations, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(99))
	if err != nil {
		return err
	}
	recomendationIds := []spotify.ID{}
	for _, song := range recomendations.Tracks {
		recomendationIds = append(recomendationIds, song.ID)
	}
	err = c.ClearRadio(ctx)
	if err != nil {
		return err
	}
	radioPlaylist, db, err := c.GetRadioPlaylist(ctx, song.Name)
	if err != nil {
		return err
	}
	_, err = db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(song.ID)))
	if err != nil {
		return err
	}
	queue := []spotify.ID{song.ID}
	for _, rec := range recomendationIds {
		exists, err := c.SongExists(db, rec)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(rec)))
			if err != nil {
				return err
			}
			queue = append(queue, rec)
		}
	}
	_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, queue...)
	if err != nil {
		return err
	}
	delay := time.Now().UnixMilli() - start
	if pos != 0 {
		pos = pos + int(delay)
	}
	err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
		PlaybackContext: &radioPlaylist.URI,
		PlaybackOffset: &spotify.PlaybackOffset{
			Position: 0,
		},
		PositionMs: pos,
	})
	if err != nil {
		if isNoActiveError(err) {
			deviceID, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
				PlaybackContext: &radioPlaylist.URI,
				PlaybackOffset: &spotify.PlaybackOffset{
					Position: 0,
				},
				DeviceID:   &deviceID,
				PositionMs: pos,
			})
			if err != nil {
				return err
			}
		}
	}
	err = c.Client().Repeat(ctx, "context")
	if err != nil {
		return err
	}
	for i := 0; i < 4; i++ {
		id := frand.Intn(len(recomendationIds)-2) + 1
		seed := spotify.Seeds{
			Tracks: []spotify.ID{recomendationIds[id]},
		}
		additional_recs, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
		if err != nil {
			return err
		}
		additionalRecsIds := []spotify.ID{}
		for _, song := range additional_recs.Tracks {
			exists, err := c.SongExists(db, song.ID)
			if err != nil {
				return err
			}
			if !exists {
				_, err = db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(song.ID)))
				if err != nil {
					return err
				}
				additionalRecsIds = append(additionalRecsIds, song.ID)
			}
		}
		_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, additionalRecsIds...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) DeleteTracksFromPlaylist(ctx *gctx.Context, tracks []spotify.ID, playlist spotify.ID) error {
	_, err := c.Client().RemoveTracksFromPlaylist(ctx, playlist, tracks...)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) SongExists(db *sql.DB, song spotify.ID) (bool, error) {
	song_id := string(song)
	sqlStmt := `SELECT id FROM radio WHERE id = ?`
	err := db.QueryRow(sqlStmt, song_id).Scan(&song_id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (c *Commands) Radio(ctx *gctx.Context) error {
	current_song, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	var seed_song spotify.SimpleTrack

	if current_song.Item != nil {
		seed_song = current_song.Item.SimpleTrack
	}
	if current_song.Item == nil {
		_, err := c.activateDevice(ctx)
		if err != nil {
			return err
		}
		tracks, err := c.Client().CurrentUsersTracks(ctx, spotify.Limit(10))
		if err != nil {
			return err
		}
		seed_song = tracks.Tracks[frand.Intn(len(tracks.Tracks))].SimpleTrack
	} else if !current_song.Playing {
		tracks, err := c.Client().CurrentUsersTracks(ctx, spotify.Limit(10))
		if err != nil {
			return err
		}
		seed_song = tracks.Tracks[frand.Intn(len(tracks.Tracks))].SimpleTrack
	}
	return c.RadioGivenSong(ctx, seed_song, current_song.Progress)
}

func (c *Commands) RefillRadio(ctx *gctx.Context) error {
	status, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	if !status.Playing {
		return nil
	}
	to_remove := []spotify.ID{}
	radioPlaylist, db, err := c.GetRadioPlaylist(ctx, "")
	if err != nil {
		return err
	}

	if status.PlaybackContext.URI != radioPlaylist.URI {
		return nil
	}

	playlistItems, err := c.Client().GetPlaylistItems(ctx, radioPlaylist.ID)
	if err != nil {
		return fmt.Errorf("orig playlist items: %w", err)
	}

	page := 0
	for {
		tracks, err := c.Client().GetPlaylistItems(ctx, radioPlaylist.ID, spotify.Limit(50), spotify.Offset(page*50))
		if err != nil {
			return fmt.Errorf("tracks: %w", err)
		}
		if len(tracks.Items) == 0 {
			break
		}
		for _, track := range tracks.Items {
			if track.Track.Track.ID == status.Item.ID {
				break
			}
			to_remove = append(to_remove, track.Track.Track.ID)
		}
		page++
	}
	if len(to_remove) > 0 {
		var trackGroups []spotify.ID
		for idx, item := range to_remove {
			if idx%100 == 0 {
				_, err = c.Client().RemoveTracksFromPlaylist(ctx, radioPlaylist.ID, trackGroups...)
				trackGroups = []spotify.ID{}
			}
			trackGroups = append(trackGroups, item)
			if err != nil {
				return fmt.Errorf("error clearing playlist: %w", err)
			}
		}
		c.Client().RemoveTracksFromPlaylist(ctx, radioPlaylist.ID, trackGroups...)
	}

	to_add := 500 - (playlistItems.Total - len(to_remove))
	playlistItems, err = c.Client().GetPlaylistItems(ctx, radioPlaylist.ID)
	if err != nil {
		return fmt.Errorf("playlist items: %w", err)
	}
	total := playlistItems.Total
	pages := int(math.Ceil(float64(total) / 50))
	randomPage := 1
	if pages > 1 {
		randomPage = frand.Intn(pages-1) + 1
	}
	playlistPage, err := c.Client().GetPlaylistItems(ctx, radioPlaylist.ID, spotify.Limit(50), spotify.Offset((randomPage-1)*50))
	if err != nil {
		return fmt.Errorf("playlist page: %w", err)
	}
	pageSongs := playlistPage.Items
	frand.Shuffle(len(pageSongs), func(i, j int) { pageSongs[i], pageSongs[j] = pageSongs[j], pageSongs[i] })
	seedCount := 5
	if len(pageSongs) < seedCount {
		seedCount = len(pageSongs)
	}
	seedIds := []spotify.ID{}
	for idx, song := range pageSongs {
		if idx >= seedCount {
			break
		}
		seedIds = append(seedIds, song.Track.Track.ID)
	}
	seed := spotify.Seeds{
		Tracks: seedIds,
	}
	recomendations, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(95))
	if err != nil {
		return err
	}
	recomendationIds := []spotify.ID{}
	for _, song := range recomendations.Tracks {
		exists, err := c.SongExists(db, song.ID)
		if err != nil {
			return fmt.Errorf("err check song existnce: %w", err)
		}
		if !exists {
			recomendationIds = append(recomendationIds, song.ID)
		}
	}
	queue := []spotify.ID{}
	for idx, rec := range recomendationIds {
		if idx > to_add {
			break
		}
		_, err = db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", rec.String()))
		if err != nil {
			return err
		}
		queue = append(queue, rec)
	}
	to_add -= len(queue)
	_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, queue...)
	if err != nil {
		return fmt.Errorf("add tracks: %w", err)
	}
	err = c.Client().Repeat(ctx, "context")
	if err != nil {
		return fmt.Errorf("repeat: %w", err)
	}
	for to_add > 0 {
		id := frand.Intn(len(recomendationIds)-2) + 1
		seed := spotify.Seeds{
			Tracks: []spotify.ID{recomendationIds[id]},
		}
		additional_recs, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
		if err != nil {
			return fmt.Errorf("get recs: %w", err)
		}
		additionalRecsIds := []spotify.ID{}
		for idx, song := range additional_recs.Tracks {
			exists, err := c.SongExists(db, song.ID)
			if err != nil {
				return fmt.Errorf("check song existence: %w", err)
			}
			if !exists {
				if idx > to_add {
					break
				}
				additionalRecsIds = append(additionalRecsIds, song.ID)
				queue = append(queue, song.ID)
			}
		}
		to_add -= len(queue)
		_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, additionalRecsIds...)
		if err != nil {
			return fmt.Errorf("add tracks to playlist: %w", err)
		}
	}
	return nil
}

func (c *Commands) ClearRadio(ctx *gctx.Context) error {
	radioPlaylist, db, err := c.GetRadioPlaylist(ctx, "")
	if err != nil {
		return err
	}
	err = c.Client().UnfollowPlaylist(ctx, radioPlaylist.ID)
	if err != nil {
		return err
	}
	db.Query("DROP TABLE IF EXISTS radio")
	configDir, _ := os.UserConfigDir()
	os.Remove(filepath.Join(configDir, "gospt/radio.json"))
	c.Client().Pause(ctx)
	return nil
}

func (c *Commands) Devices(ctx *gctx.Context) error {
	devices, err := c.Client().PlayerDevices(ctx)
	if err != nil {
		return err
	}
	return PrintDevices(devices)
}

func (c *Commands) Pause(ctx *gctx.Context) error {
	err := c.Client().Pause(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) TogglePlay(ctx *gctx.Context) error {
	current, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return c.Play(ctx)
	}
	if !current.Playing {
		return c.Play(ctx)
	}
	return c.Pause(ctx)
}

func (c *Commands) Like(ctx *gctx.Context) error {
	playing, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	err = c.Client().AddTracksToLibrary(ctx, playing.Item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Unlike(ctx *gctx.Context) error {
	playing, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	err = c.Client().RemoveTracksFromLibrary(ctx, playing.Item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Next(ctx *gctx.Context, amt int, inqueue bool) error {
	if inqueue {
		for i := 0; i < amt; i++ {
			c.Client().Next(ctx)
		}
		return nil
	}
	if amt == 1 {
		err := c.Client().Next(ctx)
		if err != nil {
			if isNoActiveError(err) {
				deviceId, err := c.activateDevice(ctx)
				if err != nil {
					return err
				}
				err = c.Client().NextOpt(ctx, &spotify.PlayOptions{
					DeviceID: &deviceId,
				})
				if err != nil {
					return err
				}
			}
			return err
		}
		return nil
	}
	// found := false
	// playingIndex := 0
	current, err := c.Client().PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return err
	}
	playbackContext := current.PlaybackContext.Type
	switch playbackContext {
	case "playlist":
		found := false
		currentTrackIndex := 0
		page := 1
		for !found {
			playlist, err := c.Client().GetPlaylistItems(ctx, spotify.ID(strings.Split(string(current.PlaybackContext.URI), ":")[2]), spotify.Limit(50), spotify.Offset((page-1)*50))
			if err != nil {
				return err
			}
			for idx, track := range playlist.Items {
				if track.Track.Track.ID == current.Item.ID {
					currentTrackIndex = idx + (50 * (page - 1))
					found = true
					break
				}
			}
			page++
		}
		c.Client().PlayOpt(ctx, &spotify.PlayOptions{
			PlaybackContext: &current.PlaybackContext.URI,
			PlaybackOffset: &spotify.PlaybackOffset{
				Position: currentTrackIndex + amt,
			},
		})
		return nil
	case "album":
		found := false
		currentTrackIndex := 0
		page := 1
		for !found {
			playlist, err := c.Client().GetAlbumTracks(ctx, spotify.ID(strings.Split(string(current.PlaybackContext.URI), ":")[2]), spotify.Limit(50), spotify.Offset((page-1)*50))
			if err != nil {
				return err
			}
			for idx, track := range playlist.Tracks {
				if track.ID == current.Item.ID {
					currentTrackIndex = idx + (50 * (page - 1))
					found = true
					break
				}
			}
			page++
		}
		c.Client().PlayOpt(ctx, &spotify.PlayOptions{
			PlaybackContext: &current.PlaybackContext.URI,
			PlaybackOffset: &spotify.PlaybackOffset{
				Position: currentTrackIndex + amt,
			},
		})
		return nil
	default:
		for i := 0; i < amt; i++ {
			c.Client().Next(ctx)
		}
	}
	return nil
}

func (c *Commands) Previous(ctx *gctx.Context) error {
	err := c.Client().Previous(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Status(ctx *gctx.Context) error {
	state, err := cache.DefaultCache().GetOrDo("state", func() (string, error) {
		state, err := c.Client().PlayerState(ctx)
		if err != nil {
			return "", err
		}
		str, err := c.FormatState(state)
		if err != nil {
			return "", nil
		}
		return str, nil
	}, 5*time.Second)
	if err != nil {
		return err
	}
	fmt.Println(state)
	return nil
}

func (c *Commands) DownloadCover(ctx *gctx.Context, args []string) error {
	destinationPath := filepath.Clean(args[0])
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(destinationPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	err = state.Item.Album.Images[0].Download(f)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Link(ctx *gctx.Context) (string, error) {
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return "", err
	}
	return state.Item.ExternalURLs["spotify"], nil
}

func (c *Commands) YoutubeLink(ctx *gctx.Context) (string, error) {
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return "", err
	}
	link := youtube.Search(state.Item.Artists[0].Name + state.Item.Name)
	return link, nil
}

func (c *Commands) LinkContext(ctx *gctx.Context) (string, error) {
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return "", err
	}
	return state.PlaybackContext.ExternalURLs["spotify"], nil
}

func (c *Commands) NowPlaying(ctx *gctx.Context, args []string) error {
	if len(args) > 0 {
		if args[0] == "force" {
			current, err := c.Client().PlayerCurrentlyPlaying(ctx)
			if err != nil {
				return err
			}
			str := FormatSong(current)
			fmt.Println(str)
			return nil
		}
	}
	song, err := cache.DefaultCache().GetOrDo("now_playing", func() (string, error) {
		current, err := c.Client().PlayerCurrentlyPlaying(ctx)
		if err != nil {
			return "", err
		}
		str := FormatSong(current)
		return str, nil
	}, 5*time.Second)
	if err != nil {
		return err
	}
	fmt.Println(song)
	return nil
}

func FormatSong(current *spotify.CurrentlyPlaying) string {
	icon := "▶"
	if !current.Playing {
		icon = "⏸"
	}
	return fmt.Sprintf("%s %s - %s", icon, current.Item.Name, current.Item.Artists[0].Name)
}

func (c *Commands) Shuffle(ctx *gctx.Context) error {
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current playstate")
	}
	err = c.Client().Shuffle(ctx, !state.ShuffleState)
	if err != nil {
		return err
	}
	ctx.Println("Shuffle set to", !state.ShuffleState)
	return nil
}

func (c *Commands) Repeat(ctx *gctx.Context) error {
	state, err := c.Client().PlayerState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current playstate")
	}
	newState := "off"
	if state.RepeatState == "off" {
		newState = "context"
	}
	// spotifyd only supports binary value for repeat, context or off, change when/if spotifyd is better
	err = c.Client().Repeat(ctx, newState)
	if err != nil {
		return err
	}
	ctx.Println("Repeat set to", newState)
	return nil
}

func (c *Commands) TrackList(ctx *gctx.Context, page int) (*spotify.SavedTrackPage, error) {
	return c.Client().CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func (c *Commands) Playlists(ctx *gctx.Context, page int) (*spotify.SimplePlaylistPage, error) {
	return c.Client().CurrentUsersPlaylists(ctx, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func (c *Commands) PlaylistTracks(ctx *gctx.Context, playlist spotify.ID, page int) (*spotify.PlaylistItemPage, error) {
	return c.Client().GetPlaylistItems(ctx, playlist, spotify.Limit(50), spotify.Offset((page-1)*50))
}

func (c *Commands) FormatState(state *spotify.PlayerState) (string, error) {
	state.Item.AvailableMarkets = []string{}
	state.Item.Album.AvailableMarkets = []string{}
	out, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return "", err
	}
	return (string(out)), nil
}

func (c *Commands) PrintPlaying(current *spotify.CurrentlyPlaying) error {
	icon := "▶"
	if !current.Playing {
		icon = "⏸"
	}
	fmt.Printf("%s %s - %s\n", icon, current.Item.Name, current.Item.Artists[0].Name)
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

func (c *Commands) SetDevice(ctx *gctx.Context, device spotify.PlayerDevice) error {
	out, err := json.MarshalIndent(device, "", " ")
	if err != nil {
		return err
	}
	configDir, _ := os.UserConfigDir()
	err = os.WriteFile(filepath.Join(configDir, "gospt/device.json"), out, 0o600)
	if err != nil {
		return err
	}
	_, err = c.activateDevice(ctx)
	if err != nil {
		return err
	}
	return nil
}

func isNoActiveError(err error) bool {
	return strings.Contains(err.Error(), "No active device found")
}

func (c *Commands) RadioFromPlaylist(ctx *gctx.Context, playlist spotify.SimplePlaylist) error {
	total := playlist.Tracks.Total
	if total == 0 {
		return fmt.Errorf("this playlist is empty")
	}
	pages := int(math.Ceil(float64(total) / 50))
	randomPage := 1
	if pages > 1 {
		randomPage = frand.Intn(pages-1) + 1
	}
	playlistPage, err := c.Client().GetPlaylistItems(ctx, playlist.ID, spotify.Limit(50), spotify.Offset((randomPage-1)*50))
	if err != nil {
		return err
	}
	pageSongs := playlistPage.Items
	frand.Shuffle(len(pageSongs), func(i, j int) { pageSongs[i], pageSongs[j] = pageSongs[j], pageSongs[i] })
	seedCount := 5
	if len(pageSongs) < seedCount {
		seedCount = len(pageSongs)
	}
	seedIds := []spotify.ID{}
	for idx, song := range pageSongs {
		if idx >= seedCount {
			break
		}
		seedIds = append(seedIds, song.Track.Track.ID)
	}
	return c.RadioGivenList(ctx, seedIds[:seedCount], playlist.Name)
}

func (c *Commands) RadioFromAlbum(ctx *gctx.Context, album spotify.SimpleAlbum) error {
	tracks, err := c.AlbumTracks(ctx, album.ID, 1)
	if err != nil {
		return err
	}
	total := tracks.Total
	if total == 0 {
		return fmt.Errorf("this playlist is empty")
	}
	pages := int(math.Ceil(float64(total) / 50))
	randomPage := 1
	if pages > 1 {
		randomPage = frand.Intn(pages-1) + 1
	}
	albumTrackPage, err := c.AlbumTracks(ctx, album.ID, randomPage)
	if err != nil {
		return err
	}
	pageSongs := albumTrackPage.Tracks
	frand.Shuffle(len(pageSongs), func(i, j int) { pageSongs[i], pageSongs[j] = pageSongs[j], pageSongs[i] })
	seedCount := 5
	if len(pageSongs) < seedCount {
		seedCount = len(pageSongs)
	}
	seedIds := []spotify.ID{}
	for idx, song := range pageSongs {
		if idx >= seedCount {
			break
		}
		seedIds = append(seedIds, song.ID)
	}
	return c.RadioGivenList(ctx, seedIds[:seedCount], album.Name)
}

func (c *Commands) RadioFromSavedTracks(ctx *gctx.Context) error {
	savedSongs, err := c.Client().CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset(0))
	if err != nil {
		return err
	}
	if savedSongs.Total == 0 {
		return fmt.Errorf("you have no saved songs")
	}
	pages := int(math.Ceil(float64(savedSongs.Total) / 50))
	randomPage := 1
	if pages > 1 {
		randomPage = frand.Intn(pages-1) + 1
	}
	trackPage, err := c.Client().CurrentUsersTracks(ctx, spotify.Limit(50), spotify.Offset(randomPage*50))
	if err != nil {
		return err
	}
	pageSongs := trackPage.Tracks
	frand.Shuffle(len(pageSongs), func(i, j int) { pageSongs[i], pageSongs[j] = pageSongs[j], pageSongs[i] })
	seedCount := 4
	seedIds := []spotify.ID{}
	for idx, song := range pageSongs {
		if idx >= seedCount {
			break
		}
		seedIds = append(seedIds, song.ID)
	}
	seedIds = append(seedIds, savedSongs.Tracks[0].ID)
	return c.RadioGivenList(ctx, seedIds, "Saved Tracks")
}

func (c *Commands) RadioGivenList(ctx *gctx.Context, song_ids []spotify.ID, name string) error {
	seed := spotify.Seeds{
		Tracks: song_ids,
	}
	recomendations, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(99))
	if err != nil {
		return err
	}
	recomendationIds := []spotify.ID{}
	for _, song := range recomendations.Tracks {
		recomendationIds = append(recomendationIds, song.ID)
	}
	err = c.ClearRadio(ctx)
	if err != nil {
		return err
	}
	radioPlaylist, db, err := c.GetRadioPlaylist(ctx, name)
	if err != nil {
		return err
	}
	queue := []spotify.ID{song_ids[0]}
	for _, rec := range recomendationIds {
		exists, err := c.SongExists(db, rec)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(rec)))
			if err != nil {
				return err
			}
			queue = append(queue, rec)
		}
	}
	_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, queue...)
	if err != nil {
		return err
	}
	err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
		PlaybackContext: &radioPlaylist.URI,
		PlaybackOffset: &spotify.PlaybackOffset{
			Position: 0,
		},
	})
	if err != nil {
		if isNoActiveError(err) {
			deviceId, err := c.activateDevice(ctx)
			if err != nil {
				return err
			}
			err = c.Client().PlayOpt(ctx, &spotify.PlayOptions{
				PlaybackContext: &radioPlaylist.URI,
				PlaybackOffset: &spotify.PlaybackOffset{
					Position: 0,
				},
				DeviceID: &deviceId,
			})
			if err != nil {
				return err
			}
		}
	}
	for i := 0; i < 4; i++ {
		id := frand.Intn(len(recomendationIds)-2) + 1
		seed := spotify.Seeds{
			Tracks: []spotify.ID{recomendationIds[id]},
		}
		additional_recs, err := c.Client().GetRecommendations(ctx, seed, &spotify.TrackAttributes{}, spotify.Limit(100))
		if err != nil {
			return err
		}
		additionalRecsIds := []spotify.ID{}
		for _, song := range additional_recs.Tracks {
			exists, err := c.SongExists(db, song.ID)
			if err != nil {
				return err
			}
			if !exists {
				_, err = db.QueryContext(ctx, fmt.Sprintf("INSERT INTO radio (id) VALUES('%s')", string(song.ID)))
				if err != nil {
					return err
				}
				additionalRecsIds = append(additionalRecsIds, song.ID)
			}
		}
		_, err = c.Client().AddTracksToPlaylist(ctx, radioPlaylist.ID, additionalRecsIds...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) activateDevice(ctx *gctx.Context) (spotify.ID, error) {
	var device *spotify.PlayerDevice
	configDir, _ := os.UserConfigDir()
	if _, err := os.Stat(filepath.Join(configDir, "gospt/device.json")); err == nil {
		deviceFile, err := os.Open(filepath.Join(configDir, "gospt/device.json"))
		if err != nil {
			return "", err
		}
		defer deviceFile.Close()
		deviceValue, err := io.ReadAll(deviceFile)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(deviceValue, &device)
		if err != nil {
			return "", err
		}
		err = c.Client().TransferPlayback(ctx, device.ID, true)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Println("YOU MUST RUN gospt setdevice FIRST")
	}
	return device.ID, nil
}

func (c *Commands) GetRadioPlaylist(ctx *gctx.Context, name string) (*spotify.FullPlaylist, *sql.DB, error) {
	configDir, _ := os.UserConfigDir()
	playlistFile, err := os.ReadFile(filepath.Join(configDir, "gospt/radio.json"))
	if errors.Is(err, os.ErrNotExist) {
		return c.CreateRadioPlaylist(ctx, name)
	}
	if err != nil {
		return nil, nil, err
	}
	var playlist *spotify.FullPlaylist
	err = json.Unmarshal(playlistFile, &playlist)
	if err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(configDir, "gospt/radio.db"))
	if err != nil {
		return nil, nil, err
	}
	return playlist, db, nil
}

func (c *Commands) CreateRadioPlaylist(ctx *gctx.Context, name string) (*spotify.FullPlaylist, *sql.DB, error) {
	// private flag doesnt work
	configDir, _ := os.UserConfigDir()
	playlist, err := c.Client().CreatePlaylistForUser(ctx, c.User(), name+" - autoradio", "Automanaged radio playlist", false, false)
	if err != nil {
		return nil, nil, err
	}
	raw, err := json.MarshalIndent(playlist, "", " ")
	if err != nil {
		return nil, nil, err
	}
	err = os.WriteFile(filepath.Join(configDir, "gospt/radio.json"), raw, 0o600)
	if err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(configDir, "gospt/radio.db"))
	if err != nil {
		return nil, nil, err
	}
	db.QueryContext(ctx, "DROP TABLE IF EXISTS radio")
	db.QueryContext(ctx, "CREATE TABLE IF NOT EXISTS radio (id string PRIMARY KEY)")
	return playlist, db, nil
}

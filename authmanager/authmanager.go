package authmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gospt/config"
	"gospt/ctx"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	auth  *spotifyauth.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func GetClient(ctx *ctx.Context) (*spotify.Client, error) {
	auth = spotifyauth.New(
		spotifyauth.WithClientID(config.Values.ClientId),
		spotifyauth.WithClientSecret(config.Values.ClientSecret),
		spotifyauth.WithRedirectURL("http://localhost:1024/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeImageUpload,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopeUserFollowModify,
			spotifyauth.ScopeUserFollowRead,
			spotifyauth.ScopeUserLibraryModify,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopeUserReadRecentlyPlayed,
			spotifyauth.ScopeUserTopRead,
			spotifyauth.ScopeStreaming,
		),
	)
	homdir, _ := os.UserHomeDir()
	if _, err := os.Stat(filepath.Join(homdir, ".config/gospt/auth.json")); err == nil {
		authFile, err := os.Open(filepath.Join(homdir, ".config/gospt/auth.json"))
		if err != nil {
			return nil, err
		}
		defer authFile.Close()
		authValue, err := ioutil.ReadAll(authFile)
		if err != nil {
			return nil, err
		}
		var tok *oauth2.Token
		err = json.Unmarshal(authValue, &tok)
		if err != nil {
			return nil, err
		}
		client := spotify.New(auth.Client(ctx, tok))
		new_token, err := client.Token()
		if err != nil {
			return nil, err
		}
		if new_token != tok {
			out, err := json.MarshalIndent(new_token, "", " ")
			if err != nil {
				panic(err.Error())
			}
			homdir, _ := os.UserHomeDir()
			err = ioutil.WriteFile(filepath.Join(homdir, ".config/gospt/auth.json"), out, 0o644)
			if err != nil {
				panic("FAILED TO SAVE AUTH")
			}
		}
		out, err := json.MarshalIndent(tok, "", " ")
		if err != nil {
			panic(err.Error())
		}
		homdir, _ := os.UserHomeDir()
		err = ioutil.WriteFile(filepath.Join(homdir, ".config/gospt/auth.json"), out, 0o644)
		if err != nil {
			panic("FAILED TO SAVE AUTH")
		}
		return client, nil
	}
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":1024", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	url := auth.AuthURL(state)
	fmt.Println("GO HERE", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("You are logged in as:", user.ID)
	return client, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	out, err := json.MarshalIndent(tok, "", " ")
	if err != nil {
		panic(err.Error())
	}
	homdir, _ := os.UserHomeDir()
	err = ioutil.WriteFile(filepath.Join(homdir, ".config/gospt/auth.json"), out, 0o644)
	if err != nil {
		panic("FAILED TO SAVE AUTH")
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

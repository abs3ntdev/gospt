package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"tuxpa.in/a/zlog/log"

	"gitea.asdf.cafe/abs3nt/gospt/src/config"
	"gitea.asdf.cafe/abs3nt/gospt/src/gctx"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	auth         *spotifyauth.Authenticator
	ch           = make(chan *spotify.Client)
	state        = "abc123"
	configDir, _ = os.UserConfigDir()
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func GetClient(ctx *gctx.Context) (*spotify.Client, error) {
	if config.Values.ClientId == "" || config.Values.ClientSecret == "" || config.Values.Port == "" {
		fmt.Println("PLEASE WRITE YOUR CONFIG FILE IN", filepath.Join(configDir, "gospt/client.yml"))
		fmt.Println("GO HERE TO AND MAKE AN APPLICATION: https://developer.spotify.com/dashboard/applications")
		fmt.Print("\nclient_id: \"idgoesherelikethis\"\nclient_secret: \"secretgoesherelikethis\"\nport:\"8888\"\n\n")
		return nil, fmt.Errorf("\nINVALID CONFIG")
	}
	auth = spotifyauth.New(
		spotifyauth.WithClientID(config.Values.ClientId),
		spotifyauth.WithClientSecret(config.Values.ClientSecret),
		spotifyauth.WithRedirectURL(fmt.Sprintf("http://localhost:%s/callback", config.Values.Port)),
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
	if _, err := os.Stat(filepath.Join(configDir, "gospt/auth.json")); err == nil {
		authFilePath := filepath.Join(configDir, "gospt/auth.json")
		authFile, err := os.Open(authFilePath)
		if err != nil {
			return nil, err
		}
		defer authFile.Close()
		tok := &oauth2.Token{}
		err = json.NewDecoder(authFile).Decode(tok)
		if err != nil {
			return nil, err
		}
		ctx.Context = context.WithValue(ctx.Context, oauth2.HTTPClient, &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				log.Trace().Interface("path", r.URL.Path).Msg("request")
				return http.DefaultTransport.RoundTrip(r)
			}),
		})
		authClient := auth.Client(ctx, tok)
		client := spotify.New(authClient)
		new_token, err := client.Token()
		if err != nil {
			return nil, err
		}
		out, err := json.MarshalIndent(new_token, "", " ")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(authFilePath, out, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to save auth")
		}
		return client, nil
	}
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", config.Values.Port), nil)
		if err != nil {
			panic(err)
		}
	}()
	url := auth.AuthURL(state)
	fmt.Println(url)
	cmd := exec.Command("xdg-open", url)
	cmd.Start()
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
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	out, err := json.MarshalIndent(tok, "", " ")
	if err != nil {
		panic(err.Error())
	}
	err = os.WriteFile(filepath.Join(configDir, "gospt/auth.json"), out, 0o644)
	if err != nil {
		panic("FAILED TO SAVE AUTH")
	}
	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

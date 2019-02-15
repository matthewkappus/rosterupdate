package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	gc "google.golang.org/api/classroom/v1"
)

// AuthHandler takes a key from the request and redirects to the Google oauth2 auth page
func (c *Classroom) AuthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: take key as post value
	key := "key"
	if err := c.Session.Set(w, r, c.Session.SID(), key); err != nil {
		panic(err)
	}
	key, err := c.Session.Get(r, c.Session.SID())
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, OAuth2LoginURL(key), http.StatusTemporaryRedirect)
}

// StartAPI is the callback handler registered to /oauth2callbackhandler
// It receives state, code parameters returned from google OA2 exchange and
// sets the token to make the api Service available
func (c *Classroom) StartAPI(w http.ResponseWriter, r *http.Request) {
	key, err := c.Session.Get(r, c.Session.SID())

	if key != r.FormValue("state") || err != nil {
		http.Error(w, "State mismatch", http.StatusInternalServerError)
		return
	}

	if err = c.initAPI(r.FormValue("code")); err != nil {
		http.Error(w, "Could not initialize Classroom api:"+err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := c.getUser()
	if err != nil {
		http.Error(w, "Could not get UserProfile:"+err.Error(), http.StatusInternalServerError)
		return
	}

	c.Session.Set(w, r, key, strings.ToLower(p.EmailAddress))
	http.Redirect(w, r, "/rosters", http.StatusTemporaryRedirect)
}

// initAPI enables the classroom token with code provided from an oauth2 token exchange
func (c *Classroom) initAPI(code string) error {
	ctx := context.Background()

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return err
	}
	client := oauth2.NewClient(ctx, config.TokenSource(ctx, token))
	c.api, err = gc.New(client)

	return err
}

// OAuth2LoginURL returns the url that directs a user to the Google Oauth access permission page
// Redirect to the url after authentication has set Session[sid]key for StartAPI to work
func OAuth2LoginURL(state string) string {
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ToJSON writes provided data to ResponseWriter and sets content-type, origin headers
func ToJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

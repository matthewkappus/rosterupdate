package handlers

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/uuid"
	"github.com/gorilla/sessions"
)

// Session provides get/set methods for accessing cookie store
type Session struct {
	store *sessions.CookieStore
	sid   string
	name  string
}

// NewSession returns session with sid and name
func NewSession(name string) *Session {
	return &Session{
		sid:   uuid.Must(uuid.NewV4()).String(),
		store: sessions.NewCookieStore([]byte(name)),
		name:  name,
	}

}

// SID returns SID created on NewSession()
func (s *Session) SID() string {
	return s.sid
}

// User returns the key stored as "user"
func (s *Session) User(r *http.Request) (email string, err error) {
	key, err := s.Get(r, s.SID())
	if err != nil {
		return "", fmt.Errorf("User not signed in: Log in with your APS email address")
	}
	email, err = s.Get(r, key)
	if err != nil {
		return "", fmt.Errorf("User not signed in: Log in with your APS email address")
	}
	return email, nil
}

// Get a string stored in session stored by key
func (s *Session) Get(r *http.Request, key string) (string, error) {
	sesh, err := s.store.Get(r, s.name)
	if err != nil {
		return "", err
	}
	v, ok := sesh.Values[key].(string)
	if !ok {
		return "", fmt.Errorf("Could not convert %v to string", v)
	}
	return v, nil
}

// Set stores key and value to provided request, response writer and errors otherwise
func (s *Session) Set(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	sesh, err := s.store.Get(r, s.name)
	if err != nil {
		return err
	}

	sesh.Values[key] = value

	return sesh.Save(r, w)
}

package session

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Manager is an implementation of the Sessioner Interface
type Manager struct {
}

// ReadCookie obtains the session token from the requests cookies, which come with every request
func (s *Manager) ReadCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "cookie is not set", errors.Wrap(err, "error: cookie is not set")
		}
		// For any other type of error, return a bad request status
		return "bad req", errors.Wrap(err, "error: getting the cookie")
	}
	st := c.Value

	return st, nil
}

//SetCookie  put the client cookie for "session_token" as the session token
func (s *Manager) SetCookie(w http.ResponseWriter, sessionToken string) {

	// set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

// NewST  retuns new  uniq session token
func (s *Manager) NewST() string {
	st := uuid.NewV4().String()
	return st
}

// SM is instance of SessionManager stucture
var SM = &Manager{}
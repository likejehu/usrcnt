package helpers

import (
	"net/http"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// SessionManager is an implementation of the Sessioner Interface
type SessionManager struct {
	ST string
}

// ReadCookie obtains the session token from the requests cookies, which come with every request
func (s *SessionManager) ReadCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return "cookie is not set", errors.Wrap(err, "error: cookie is not set")
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return "bad req", errors.Wrap(err, "error: getting the cookie")
	}
	s.ST = c.Value
	return s.ST, nil
}

// NewST  retuns new  uniq session token
func (s *SessionManager) NewST() string {
	s.ST = uuid.NewV4().String()
	return s.ST
}

// SM is instance of SessionManager stucture
var SM = &SessionManager{}

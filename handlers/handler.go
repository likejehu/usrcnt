package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

//Storer is  interface for  basic Key/Value (real and mock) datastorage for links
type Storer interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

//Sessioner is  interface for sessions manager
type Sessioner interface {
	ReadCookie(w http.ResponseWriter, r *http.Request) (string, error)
	NewST() string
}

// Handler is struct for handlers
type Handler struct {
	Cache   Storer
	Session Sessioner
}

//Hello is handler that creates new session and deals with logic
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) (err error) {
	usrCountkey := "usrcountkey"
	sessionToken := ""
	sessionToken, err = h.Session.ReadCookie(w, r)
	if sessionToken == "bad req" {
		log.Print(err)
		return err
	}
	if sessionToken == "cookie is not set" {
		// Create a new random session token with uuid
		sessionToken = h.Session.NewST()
		// Set the token in the cache
		// The token has an expiry time of 120 seconds
		_, err = h.Cache.Do("SETEX", sessionToken, "120")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.Wrap(err, "error: getting the result with SETEX"))
			return err

		}
	}

	//Set the client cookie for "session_token" as the session token
	// set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})

	// get the token of the user from our cache
	response, err := h.Cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return errors.Wrap(err, "error: getting the session token")
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return errors.Wrap(err, "error: token is not present in cache")
	}
	usrCountVal, err := h.Cache.Do("INCR", usrCountkey)
	w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
	w.Write([]byte(fmt.Sprintf("usercount: %v", usrCountVal)))
	return err
}

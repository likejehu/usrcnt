package handlers

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

//Storer is  interface for  basic Key/Value (real and mock) datastorage for links
type Storer interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

//Sessioner is  interface for sessions manager
type Sessioner interface {
	ReadCookie(w http.ResponseWriter, r *http.Request) (string, error)
	SetCookie(w http.ResponseWriter, sessionToken string)
	NewST() string
}

// Handler is struct for handlers
type Handler struct {
	Cache   Storer
	Session Sessioner
}

//Hello is handler that creates new session and deals with logic
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) (err error) {
	usrCountKey := "usrcountkey"
	sessionToken := ""
	sessionToken, err = h.Session.ReadCookie(w, r)
	if sessionToken == "bad req" {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
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
		h.Session.SetCookie(w, sessionToken)
		res, err := h.Cache.Do("INCR", usrCountKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.Wrap(err, "error: setting the result with INCR"))
			return err

		}
		usrCountVal, ok := res.(string)
		if ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(usrCountVal))
		}
	}
	res, err := h.Cache.Do("GET", usrCountKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errors.Wrap(err, "error: getting the result with GET"))
		return err

	}
	usrCountVal, ok := res.(string)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(usrCountVal))
	}
	log.Print("This is the end / Beautiful friend")
	return nil
}

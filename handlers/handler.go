package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
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

const usrCountKey = "usrcountkey"

//Hello is handler that creates new session and deals with logic
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var usrCountVal int
	sessionToken, err := h.Session.ReadCookie(w, r)
	log.Print("session token is ", sessionToken)
	if sessionToken == "bad req" {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)

	}
	if sessionToken == "cookie is not set" {
		// Create a new random session token with uuid
		sessionToken = h.Session.NewST()
		// Set the token in the cache
		// The token has an expiry time of 120 seconds
		_, err = h.Cache.Do("SETEX", sessionToken, "120", sessionToken)
		log.Print("session token is ", sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.Wrap(err, "error: settin  with SETEX"))

		}
		h.Session.SetCookie(w, sessionToken)
		// if usrCountKey does not exist set it value to zero
		_, err = h.Cache.Do("SETNX", usrCountKey, 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.Wrap(err, "error: settin  with SETNX"))

		}
		res, err := redis.Int(h.Cache.Do("INCR", usrCountKey))
		log.Print("after INCR usrCountVal is now: ", res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.Wrap(err, "error: settin with INCR"))
		}
	} else {
		usrCountVal, err = redis.Int(h.Cache.Do("GET", usrCountKey))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(errors.Wrap(err, "error: getting the result with GET"))
		}
		s := strconv.Itoa(usrCountVal)
		log.Print("usrCountVal is ", usrCountVal)
		log.Print("This is the end / Beautiful friend")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s))
	}
}

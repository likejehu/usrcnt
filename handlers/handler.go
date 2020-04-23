package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/likejehu/usrcnt/session"
	"github.com/pkg/errors"
)

//Store is  interface for  basic Key/Value (real and mock) datastorage
type Store interface {
	Set(key string, value string) error
	Get(key string) (int, error)
	SETNXToZero(key string) error
	Increment(key string) (int, error)
	Exists(key string) (int, error)
}

//SessionManager is  interface for sessions manager
type SessionManager interface {
	ReadCookie(w http.ResponseWriter, r *http.Request) (string, error)
	SetCookie(w http.ResponseWriter, sessionToken string)
	NewST() string
}

// Handler is struct for handlers
type Handler struct {
	Cache   Store
	Session SessionManager
}

const usrCountKey = "usrcountkey"

//Hello is handler that creates new session and deals with logic
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	var usrCountVal int
	sessionToken, err := h.Session.ReadCookie(w, r)
	log.Print("session token is ", sessionToken)
	log.Print("err is ", err)
	if err != nil {
		if err == session.ErrorNotSet {
			// Create a new random session token with uuid
			sessionToken = h.Session.NewST()
			log.Print("session token is ", sessionToken)
			// Set the token in the cache
			// The token has an expiry time of 120 seconds
			err = h.Cache.Set(sessionToken, sessionToken)
			log.Print("session token is ", sessionToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}
			h.Session.SetCookie(w, sessionToken)
			// if usrCountKey does not exist set it value to zero
			err = h.Cache.SETNXToZero(usrCountKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}
			usrCountVal, err := h.Cache.Increment(usrCountKey)
			log.Print("after INCR usrCountVal is now: ", usrCountVal)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return
			}
		} else {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		//check is session token exists in cache
		e, err := h.Cache.Exists(sessionToken)
		log.Print("exists: ", e)
		if e == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad token!"))
			log.Print(err)
			log.Print("session token:" + sessionToken + " does not exist")
			return
		}
	}
	usrCountVal, err = h.Cache.Get(usrCountKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errors.Wrap(err, "error: getting the result with GET"))
		return
	}
	s := strconv.Itoa(usrCountVal)
	log.Print("usrCountVal is ", s)
	log.Print("This is the end / Beautiful friend")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

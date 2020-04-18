package handlers

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

//Storer is  interface for  basic Key/Value (real and mock) datastorage for links
type Storer interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

// Handler is struct for handlers
type Handler struct {
	Cache Storer
}

//Hello is handler that creates new session and deals with logic
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	usrCountkey := "usrcountkey"
	// Create a new random session token with uuid
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err := h.Cache.Do("SETEX", sessionToken, "120")
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Set the client cookie for "session_token" as the session token
	// set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})

	// obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken = c.Value

	// get the name of the user from our cache, where we set the session token
	response, err := h.Cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	usrCountVal, err := h.Cache.Do("INCR", usrCountkey)
	w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
	w.Write([]byte(fmt.Sprintf("usercount: %v", usrCountVal)))
}

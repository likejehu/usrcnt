package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/likejehu/usrcnt/db"
	"github.com/likejehu/usrcnt/handlers"
	"github.com/likejehu/usrcnt/session"
)

func main() {
	// cache is instance of redis storage
	var cache = *db.NewRedisStore("redis://localhost:6379")
	defer cache.Close()
	handler := handlers.Handler{
		Cache:   cache,
		Session: session.SM,
	}
	router := httprouter.New()
	router.GET("/hello", handler.Hello)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}

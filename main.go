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
	// conn is instance of redis.Conn
	var conn = db.RedisCache.Client
	defer conn.Close()
	handler := handlers.Handler{
		Cache:   db.RedisCache,
		Session: session.SM,
	}
	router := httprouter.New()
	router.GET("/hello", handler.Hello)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}

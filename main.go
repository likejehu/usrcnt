package main

import (
	"log"
	"net/http"

	"github.com/likejehu/usrcnt/db"
	"github.com/likejehu/usrcnt/handlers"
)

func main() {
	// cache is instance of redis storage
	var cache = *db.NewRedisStore("redis://localhost:6379")
	defer cache.Close()
	handler := handlers.Handler{
		Cache: cache,
	}
	http.HandleFunc("/hello", handler.Hello)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

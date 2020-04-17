package main

import (
	"log"
	"net/http"

	"github.com/likejehu/usrcnt/db"
	"github.com/likejehu/usrcnt/handlers"
)

func main() {
	handler := handlers.Handler{
		Cache: db.Cache,
	}
	http.HandleFunc("/hello", handler.Hello)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

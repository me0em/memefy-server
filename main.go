package main

import (
	"log"
	"memefy/api"
	"net/http"
)


func runServer() {
	// TODO: (GET /) func handler, HTTP/2 pushed css
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/user", api.CreateUser)
	http.HandleFunc("/api/memes", api.GetMemes)
	http.HandleFunc("/api/reaction", api.HandleReaction)
	http.HandleFunc("/api/test", api.TestThings)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	runServer()
}

package main

import (
	"log"
	"memefy-server/api"
	"net/http"
)


func runServer() {
	// TODO: (GET /) func handler, HTTP/2 pushed css
	// Configs etc
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/user", api.CreateUser)
	http.HandleFunc("/api/memes", api.ThrowMemes)
	http.HandleFunc("/api/reaction", api.HandleReaction)
	http.HandleFunc("/api/test", api.TestThings)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	runServer()
}

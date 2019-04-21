package main

import (
	"fmt"
	"log"
	"net/http"

	"memefy-server/api"
	"memefy-server/config"
)

func runServer() {
	// TODO: (GET /) func handler, HTTP/2 pushed css
	// Configs etc
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/user", api.CreateUser)
	http.HandleFunc("/api/memes", api.SendMemes)
	http.HandleFunc("/api/reaction", api.GetReaction)
	http.HandleFunc("/api/refresh", api.RefreshToken)
	http.HandleFunc("/api/test", api.TestThings)
	log.Fatal(http.ListenAndServe(config.ServerPort, nil))
}

func main() {
	fmt.Println("Run server...")
	api.InitDB(config.DBInitReq)
	fmt.Println("Init the database successfully...")
	runServer()
}

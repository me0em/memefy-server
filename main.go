package main

import (
	"fmt"
	"log"
	"memefy-server/api"
	"net/http"
)

func runServer() {
	// TODO: (GET /) func handler, HTTP/2 pushed css
	// Configs etc
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/user", api.CreateUser)
	http.HandleFunc("/api/memes", api.SendMemes)
	http.HandleFunc("/api/reaction", api.GetReaction)
	http.HandleFunc("/api/test", api.TestThings)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func main() {
	fmt.Println("Run server...")
	api.InitDB("root@tcp(127.0.0.1:3306)/memefy?parseTime=true")
	fmt.Println("Init the database successfully...")
	runServer()
}

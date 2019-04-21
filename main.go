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
	http.HandleFunc("/api/refresh", api.RefreshToken)
	http.HandleFunc("/api/test", api.TestThings)
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func main() {
	fmt.Println("Run server...")
	//api.InitDB("root@tcp(127.0.0.1:3	306)/memefy?parseTime=true")
	api.InitDB("http://default::tyz1214@127.0.0.1:8123/memefy")
	fmt.Println("Init the database successfully...")
	runServer()
}

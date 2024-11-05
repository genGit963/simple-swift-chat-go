package main

import (
	"log"
	"net/http"

	"github.com/genGit963/simple-swift-chat-go/internal/handlers"
)

func main() {

	mux := routes()

	log.Println("Starting the channel listener")
	go handlers.ListenToWebsocketChannel()

	log.Println("Server on 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err.Error())
	}
}

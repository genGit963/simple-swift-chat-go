package main

import (
	"log"
	"net/http"

	"github.com/genGit963/simple-swift-chat-go/internal/handlers"
	"github.com/genGit963/simple-swift-chat-go/utils/errorutils"
)

func main() {

	mux := routes()

	log.Println("Starting the channel listener")
	go handlers.ListenToWebsocketChannel()

	log.Println("Server on 8080")
	err := http.ListenAndServe(":8080", mux)
	errorutils.AnyErrorCaptureLog("Server listening", 0, err)
}

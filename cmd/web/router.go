package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/genGit963/simple-swift-chat-go/internal/handlers"
)

func routes() http.Handler {

	mux := pat.New()

	// route home
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WSEndpoint))

	return mux
}

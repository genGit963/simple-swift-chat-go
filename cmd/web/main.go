package main

import (
	"log"
	"net/http"
)

func main() {

	mux := routes()

	log.Println("Server on 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println(err.Error())
	}
}

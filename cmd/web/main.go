package main

import (
	"fmt"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func greetHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I m greeting for go lang, Hello at %s", time.Now())
}

func main() {
	http.HandleFunc("/da", greet)
	http.HandleFunc("/hello", greetHello)
	http.ListenAndServe(":8080", nil)
}

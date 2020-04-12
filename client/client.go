package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("client")))
	log.Println("Serving SSE test client on localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

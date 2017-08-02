package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/request", racer)
	http.HandleFunc("/solution", solution)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/request", racer)
	http.HandleFunc("/solution", solution)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// algo.Process("Mike Tyson", "International Boxing Federation", req, path, 10000, &jobqueue.Sem)

}

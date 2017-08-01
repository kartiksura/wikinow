package main

import (
	"log"
	"time"

	"github.com/kartiksura/wikinow/algo"
	"github.com/kartiksura/wikinow/jobqueue"
)

func main() {
	req, err := algo.SetNewRequest(50000)
	if err != nil {
		log.Println(err)
		return
	}
	var path []string
	algo.Process("Mike Tyson", "International Boxing Federation", req, path, 10000, &jobqueue.Sem)

	duration := time.Minute
	time.Sleep(duration)
}

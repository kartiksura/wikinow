package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kartiksura/wikinow/algo"
	"github.com/kartiksura/wikinow/jobqueue"
)

//Request represents a job request
type Request struct {
	ID string
}

//JobStatus represents the status of the job
type JobStatus struct {
	ID             string
	Status         string
	Path           []string
	ProcessingTime string
}

func racer(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	src := queryValues.Get("src")
	dst := queryValues.Get("dst")
	to := queryValues.Get("timeOut")

	if src == "" || dst == "" || to == "" {
		http.Error(w, "Mandatory params missing(src/dst/timeOut)", 200)
	}
	log.Println("New request:", src, dst, to)

	j, err := algo.NewJob(src, dst, to)
	if err != nil {
		log.Println(err)
		return
	}
	request := Request{ID: j.ReqID}
	js, err := json.Marshal(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	err = algo.SetJob(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go algo.Process(j, &jobqueue.Sem)
}

func solution(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	id := queryValues.Get("id")
	if id == "" {
		http.Error(w, "Mandatory params missing(id)", 200)

	}
	j := JobStatus{ID: id}
	if algo.CheckIfJobAlive(id) == false {
		j.Status = "TIMEOUT"
	}
	if algo.CheckIfJobSolved(id) == true {
		job, err := algo.GetJob(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ans, err := algo.GetSolution(job.Src, job.Dst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		j.Path = ans
		j.Status = job.State
		j.ProcessingTime = job.Latency.Sub(job.ReqTime).String()
	}
	js, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

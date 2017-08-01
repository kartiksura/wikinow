package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kartiksura/wikinow/algo"
	"github.com/kartiksura/wikinow/jobqueue"
)

type Request struct {
	ID string
}
type JobStatus struct {
	ID     string
	Status string
	Path   []string
}

func racer(w http.ResponseWriter, r *http.Request) {
	log.Println("new job request")
	queryValues := r.URL.Query()
	src := queryValues.Get("src")
	dst := queryValues.Get("dst")
	to := queryValues.Get("to")

	if src == "" || dst == "" || to == "" {
		http.Error(w, "Mandatory params missing(src/dst/to)", 200)
	}
	log.Println("New request:", src, dst, to)

	timeOut, err := strconv.Atoi(to)
	if err != nil {
		http.Error(w, err.Error(), 200)
	}

	req, err := algo.SetNewRequest(src, dst, timeOut)
	if err != nil {
		log.Println(err)
		return
	}
	request := Request{ID: req}
	js, err := json.Marshal(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	go algo.Process(src, dst, req, nil, timeOut, &jobqueue.Sem)
}

func solution(w http.ResponseWriter, r *http.Request) {
	log.Println("new solution request")
	queryValues := r.URL.Query()
	id := queryValues.Get("id")
	if id == "" {
		http.Error(w, "Mandatory params missing(id)", 200)

	}
	j := JobStatus{ID: id}
	if algo.CheckIfReqAlive(id) == false {
		j.Status = "TimeOut"
	}
	if algo.CheckIfReqSolved(id) == true {
		src, dst, err := algo.GetRequestDetail(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ans, err := algo.Solution(src, dst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		soln := strings.Split(ans, "|")
		j.Path = soln
		j.Status = "Solved"
	}
	js, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {

	http.HandleFunc("/request", racer)
	http.HandleFunc("/solution", solution)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// algo.Process("Mike Tyson", "International Boxing Federation", req, path, 10000, &jobqueue.Sem)

}

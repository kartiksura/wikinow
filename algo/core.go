package algo

import (
	"log"
	"strings"

	"github.com/kartiksura/wikinow/scraper"
)

//Process is the core functionality
func Process(j Job, sem *chan bool) {
	defer func() { <-*sem }()
	log.Println("Processing Job:", j)

	if CheckIfJobAlive(j.ReqID) == false || CheckIfJobSolved(j.ReqID) == true {
		log.Println("Returning request prematurely")
		return
	}

	for i := range j.Path {
		if j.Path[i] == j.Src {
			log.Println("Found loop")
			return
		}
	}
	ans, err := GetSolution(j.Src, j.Dst)
	if err == nil {
		log.Print("solution found from cache")
		SetJobState(j.ReqID, "SOLVED_FROM_CACHE")
		SetSolution(append(j.Path, ans...))
		return
	}
	j.Path = append(j.Path, j.Src)

	//calling the scraper
	links, err := scraper.ProcessTitle(j.Src)
	for i := range links {
		if CheckIfJobAlive(j.ReqID) == false || CheckIfJobSolved(j.ReqID) == true {
			log.Println("Returning request prematurely")
			return
		}
		if links[i][0] == j.Dst {
			//solution found
			SetJobState(j.ReqID, "SOLVED")
			SetSolution(append(j.Path, j.Dst))
		}

		if strings.HasPrefix(links[i][0], "File") == false && strings.HasPrefix(links[i][0], "Category") == false {
			SetSolution(append(j.Path, links[i][0]))
			j.Src = links[i][0]
			j.Path = append(j.Path, j.Src)
			EnqueueJob(j)
		}

	}
}

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
	//	path = append(path, src)
	ans, err := GetSolution(j.Src, j.Dst)
	if err != nil {
		log.Print("solution notfound")
	} else {
		log.Print("solution found")

		SetJobState(j.ReqID, "SOLVED")
		SetSolution(append(j.Path, ans...))
		// log.Fatal("done1")

		return
	}
	log.Println("calling scraper")
	j.Path = append(j.Path, j.Src)

	links, err := scraper.ProcessTitle(j.Src)
	// log.Fatal("done0")

	for i := range links {
		log.Println("checking link: ", links[i][0])
		if CheckIfJobAlive(j.ReqID) == false || CheckIfJobSolved(j.ReqID) == true {
			log.Println("Returning request prematurely")
			return
		}
		if links[i][0] == j.Dst {
			log.Print("FINALLLLLLLLLLLLLLLLLLLLLL******88*******************************************************************************************************************************************************************************************************************************************************************LLLLLLLLLLLLLLLLLLLLLLLLLL")
			//solution found
			SetJobState(j.ReqID, "SOLVED")
			SetSolution(append(j.Path, j.Dst))
			// log.Fatal("done2")
		}

		if strings.HasPrefix(links[i][0], "File") == false && strings.HasPrefix(links[i][0], "Category") == false {
			log.Println("processing link: ", links[i][0])

			SetSolution(append(j.Path, links[i][0]))
			log.Println("SETTING MID SOLUTION:", append(j.Path, links[i][0]))
			// Process(links[i][0], dst, reqID, append(path, src), to, sem)
			j.Src = links[i][0]
			j.Path = append(j.Path, j.Src)
			EnqueueJob(j)
		}

	}
}

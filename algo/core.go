package algo

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/kartiksura/wikinow/cache"
	"github.com/kartiksura/wikinow/scraper"
)

type Job struct {
	Src   string
	Dst   string
	ReqID string
	Path  []string
}

//Process is the core functionality
func Process(src, dst string, reqID string, path []string, to int, sem *chan bool) {
	defer func() { <-*sem }()
	log.Println("src:", src, "dst:", dst)

	if CheckIfReqAlive(reqID) == false || CheckIfReqSolved(reqID) == true {
		log.Println("returning request prematurely")

		return
	}

	src = strings.TrimSpace(src)
	dst = strings.TrimSpace(dst)
	for i := range path {
		if path[i] == src {
			log.Println("Found loop")
			return
		}
	}
	//	path = append(path, src)
	ans, err := Solution(src, dst)
	if err != nil {
		log.Print("solution notfound")
	} else {
		log.Print("solution found")

		SetState(reqID)
		SetSolution(append(path, ans))
		return
	}
	log.Println("calling scraper")

	links, err := scraper.ProcessTitle(src)
	for i := range links {
		log.Println("checking link: ", links[i][0])

		if links[i][0] == dst {
			log.Print("FINALLLLLLLLLLLLLLLLLLLLLL****************************************************************************LLLLLLLLLLLLLLLLLLLLLLLLLL")
			//solution found
			SetState(reqID)
			path = append(path, src)
			SetSolution(append(path, dst))
		}

		if strings.HasPrefix(links[i][0], "File") == false && strings.HasPrefix(links[i][0], "Category") == false {
			log.Println("processing link: ", links[i][0])

			SetSolution(append(path, links[i][0]))
			// _, err := Solution(src, links[i][0])
			// if err == nil {
			// 	log.Print("mid solution found", src, links[i][0])
			// 	continue
			// }
			log.Println("SETTING SOLN:", append(path, links[i][0]))
			// Process(links[i][0], dst, reqID, append(path, src), to, sem)
			j := Job{Src: links[i][0], Dst: dst, ReqID: reqID, Path: append(path, src)}
			bolB, _ := json.Marshal(j)
			cache.EnQueue("JOBS", string(bolB))
		}

	}
}

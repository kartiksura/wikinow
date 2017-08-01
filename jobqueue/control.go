package jobqueue

import (
	"encoding/json"
	"log"

	"github.com/kartiksura/wikinow/algo"
	"github.com/kartiksura/wikinow/cache"
)

var Sem chan bool

func init() {
	Sem = make(chan bool, 5000)
	go Dispatcher()
}
func Dispatcher() {

	for {
		Sem <- true

		job, err := cache.DeQueue("JOBS")
		if err == nil {
			log.Println("Jobs dequed:", job)

			var j algo.Job
			json.Unmarshal([]byte(job), &j)
			go algo.Process(j.Src, j.Dst, j.ReqID, j.Path, 100, &Sem)
		}

	}

}

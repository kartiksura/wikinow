package jobqueue

import (
	"log"

	"github.com/kartiksura/wikinow/algo"
)

var Sem chan bool

func init() {
	Sem = make(chan bool, 5000)
	go Dispatcher()
}

//Dispatcher pulls the job from the redis and maintains the concurrency of the no of jobs running
func Dispatcher() {
	for {
		Sem <- true

		job, err := algo.DequeueJob()
		if err == nil {
			log.Printf("Jobs dequed:%+v\n\n\n\n\n\n\n\n\n\n", job)

			go algo.Process(job, &Sem)
		}

	}

}

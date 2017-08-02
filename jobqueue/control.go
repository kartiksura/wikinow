package jobqueue

import (
	"time"

	"github.com/kartiksura/wikinow/algo"
	"github.com/kartiksura/wikinow/conf"
)

var Sem chan bool

func init() {
	Sem = make(chan bool, conf.WikiNowConfig.Server.ConcurrentJobs)
	go Dispatcher()
}

var delay = 1

//Dispatcher pulls the job from the redis and maintains the concurrency of the no of jobs running
func Dispatcher() {
	for {
		Sem <- true

		job, err := algo.DequeueJob()
		if err == nil {
			go algo.Process(job, &Sem)
		} else {
			//exponential delay
			delay = delay * 2
			if delay > 10 {
				delay = 1
			}
			duration := time.Duration(delay) * time.Second // Pause for 10 seconds
			time.Sleep(duration)

		}

	}

}

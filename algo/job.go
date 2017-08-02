package algo

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kartiksura/wikinow/cache"
)

//Job has the details of the job with ReqID
type Job struct {
	ReqID   string
	Src     string
	Dst     string
	ReqTime time.Time
	TTL     int64
	Path    []string
	State   string
	Latency time.Time
}

//NewJob creates a new job
func NewJob(src, dst string, to string) (Job, error) {
	j := Job{}
	j.Src = strings.TrimSpace(src)
	j.Dst = strings.TrimSpace(dst)
	j.ReqTime = time.Now()

	timeOut, err := strconv.Atoi(to)
	if err != nil {
		return j, err
	}
	j.TTL = int64(timeOut)

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		return j, err
	}

	j.ReqID = strings.Trim(string(uuid), "\n")
	log.Println("New Job:", j)
	return j, nil
}

//GetJob returns the value for the key in the redis
func GetJob(id string) (Job, error) {
	ans, err := cache.GetString("REQID:" + id)
	var j Job
	err = json.Unmarshal([]byte(ans), &j)
	log.Println("job fetched", j)
	return j, err
}

//SetJob stores the value against the key in redis
func SetJob(j Job) error {

	strB, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return cache.SetString("REQID:"+j.ReqID, string(strB))

}

//DequeueJob returns the value for the key in the redis
func DequeueJob() (Job, error) {
	var j Job
	ans, err := cache.DeQueue("JOBS")
	if err != nil {
		return j, err
	}

	err = json.Unmarshal((ans).([]byte), &j)
	return j, err
}

//EnqueueJob stores the value against the key in redis
func EnqueueJob(j Job) error {

	strB, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return cache.EnQueue("JOBS", string(strB))

}

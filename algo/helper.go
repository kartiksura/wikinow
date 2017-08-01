package algo

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kartiksura/wikinow/cache"
)

func CheckIfReqAlive(reqID string) bool {
	ans, err := cache.Get("REQID_TO:" + reqID)
	log.Println("ReqID:", reqID, ":", ans)
	if err != nil {
		return false
	}

	i, err := strconv.ParseInt(ans, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	now := time.Now()
	log.Println("Now:", now.Unix())
	if tm.After(now) {
		log.Println("Req live:", reqID)
		return true
	}
	return false
}

func SetNewRequest(to int64) (s string, err error) {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		return
	}
	s = string(uuid)
	log.Println("New job:", s)
	now := time.Now()
	err = cache.Set("REQID_TO:"+s, fmt.Sprint(now.Unix()+to))
	if err != nil {
		return
	}
	return s, nil
}

func Solution(src, dst string) (string, error) {
	ans, err := cache.Get("SOLUTION:" + src + ":" + dst)
	if err != nil {
		return "", err
	}
	return ans, nil
}

func KillReq(reqID string) error {
	return cache.Set("REQID_STATUS:"+reqID, "ERROR")
}

func SetSolution(path []string) error {
	log.Println("SETTTING SOLUTION:", path)
	for i := 0; i < len(path)-1; i++ {
		log.Println("SOLUTION:"+path[i]+":"+path[len(path)-1], strings.Join(path[i:], "|"))
		err := cache.Set("SOLUTION:"+path[i]+":"+path[len(path)-1], strings.Join(path[i:], "|"))
		if err != nil {
			return err
		}
	}
	return nil
}

func SetState(reqID string) error {

	return cache.Set("REQID_STATUS:"+reqID, "SOLVED")
}

func CheckIfReqSolved(reqID string) bool {
	ans, err := cache.Get("REQID_STATUS:" + reqID)
	if err != nil {
		return false
	}

	if ans == "SOLVED" {
		return true
	}
	return false
}

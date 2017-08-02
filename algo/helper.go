package algo

import (
	"fmt"
	"strings"
	"time"

	"github.com/kartiksura/wikinow/cache"
)

//CheckIfJobAlive returns true if the ttl is valid
func CheckIfJobAlive(reqID string) bool {
	ans, err := GetJob(reqID)
	if err != nil {
		return false
	}
	req := ans.ReqTime
	ttl := ans.TTL
	if req.Add(time.Duration(ttl)*time.Second).After(time.Now()) == true {
		return true
	}
	return false
}

//GetSolution returns the existing solution for the path
func GetSolution(src, dst string) ([]string, error) {
	ans, err := cache.GetString("SOLUTION:" + src + ":" + dst)
	if err != nil {
		return nil, fmt.Errorf("solution not found: %v", err)
	}
	return strings.Split(ans, "|"), nil
}

//SetSolution sets the path for future re-use
func SetSolution(path []string) error {
	for i := 0; i < len(path)-1; i++ {
		err := cache.SetString("SOLUTION:"+path[i]+":"+path[len(path)-1], strings.Join(path[i:], "|"))
		if err != nil {
			return err
		}
	}
	return nil
}

//SetJobState sets the state of the job
func SetJobState(reqID string, state string) error {
	j, err := GetJob(reqID)
	if err != nil {
		return err
	}
	j.State = state
	j.Latency = time.Now()
	return SetJob(j)
}

//CheckIfJobSolved checks if we need to still work on the req
func CheckIfJobSolved(reqID string) bool {
	j, err := GetJob(reqID)
	if err != nil {
		return false
	}
	if j.State == "SOLVED" || j.State == "SOLVED_FROM_CACHE" {
		return true
	}
	return false
}

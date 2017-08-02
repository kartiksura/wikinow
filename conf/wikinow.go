package conf

import (
	"log"

	"gopkg.in/gcfg.v1"
)

const (
	gconfFile = "wikinow.conf"
)

//WikiNowConfig is the config for setting the params
var WikiNowConfig struct {
	Redis struct {
		HostPort string
	}
	Server struct {
		ConcurrentJobs int
	}
}

func init() {
	err := gcfg.ReadFileInto(&WikiNowConfig, gconfFile)
	if err != nil {
		log.Fatalf("conf: error: conf.init: %s", err.Error())
	}
}

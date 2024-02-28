package config

import "strings"

type Configuration struct {
	Dir         string
	RDBFilename string
}

var (
	dir         = "DIR"
	rdbfilename = "RDBFILENAME"
)

var Config = &Configuration{}

func (config Configuration) Get(key string) string {
	switch strings.ToUpper(key) {
	case dir:
		return config.Dir
	case rdbfilename:
		return config.RDBFilename
	default:
		return ""
	}
}

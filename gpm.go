package gpm

import (
	"github.com/BurntSushi/toml"

	// "log"
)

type Config struct {
	Cron  []*Process
	Event []*Process
	Boot  []*Process
}

type Event struct {
	Name string
	Data map[string]interface{}
}

type Status struct {
	Processes map[string]*Process
}

func ParseFile(file string) (*Config, error) {

	procs := &Config{}

	_, err := toml.DecodeFile(file, procs)

	return procs, err
}

// func ()

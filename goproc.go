package goproc

import (
	"github.com/BurntSushi/toml"

	// "log"
)

type Config struct {
	Process       []*ProcessTemplate
	MaxConcurrent int
}

type MonitorMessage struct {
	Process  *Process
	ExitCode int
}

type Event struct {
	Name string
	Data *map[string]interface{}
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

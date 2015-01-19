package goproc

import (
	"github.com/BurntSushi/toml"

	// "log"
)

type ProcessHook func(action string, proc *Process)
type LogFilter func(logfile string, proc *Process) string

type Config struct {
	Process       []*ProcessTemplate
	MaxConcurrent int

	LogFilters   []LogFilter
	ProcessHooks []ProcessHook
}

type MonitorMessage struct {
	Process  *Process
	ExitCode int
}

type Trigger struct {
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

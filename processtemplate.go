package goproc

import (
	// "labix.org/v2/mgo/bson"
	"bytes"
	"encoding/json"
	// "fmt"
)

type ProcessTemplate struct {
	Args         []string
	Command      string
	Cron         string
	Cwd          string
	ErrFile      string
	Event        string
	KeepAlive    bool
	LogFile      string
	Name         string
	RespawnLimit int
}

func (t *ProcessTemplate) NewProcess() *Process {
	proc := &Process{
		Template: t,
		Args:     t.Args,
	}

	return proc
}

func concatStrings(strs ...string) string {
	var buffer bytes.Buffer

	for _, str := range strs {
		buffer.WriteString(str)
	}
	return buffer.String()
}

func (t *ProcessTemplate) NewProcessWithEvent(evt *Event) *Process {
	// Replace an arg that's :event with the json representation of the event data

	args := make([]string, 0)

	for _, arg := range t.Args {
		if arg == ":json" {
			marshaled, _ := json.Marshal(evt.Data)
			args = append(args, string(marshaled))
		} else if arg == ":flags" {
			// Only pass through flags if they're strings
			for k, v := range *evt.Data {
				marshaled, _ := json.Marshal(v)
				args = append(args, concatStrings("--", k, " ", string(marshaled)))
			}
		} else {
			args = append(args, arg)
		}
	}

	proc := t.NewProcess()
	proc.Args = args
	return proc
}

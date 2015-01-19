package goproc

import (
	// "labix.org/v2/mgo/bson"
	"bytes"
	"encoding/json"
	"github.com/robfig/cron"
	"github.com/twinj/uuid"
	// "fmt"
)

type ProcessTemplate struct {
	Args         []string
	Command      string
	Cron         string
	Cwd          string
	KeepAlive    bool
	LogFile      string
	Name         string
	RespawnLimit int
	AutoStart    bool
	RunCount     int
	ResetLog     bool
	manager      *Manager
}

func (t *ProcessTemplate) NewProcess() *Process {
	proc := &Process{
		Template: t,
		Args:     t.Args,
	}

	id := uuid.NewV1()
	proc.Uuid = id.String()

	return proc
}

func (t *ProcessTemplate) RegisterCron(c *cron.Cron) {
	c.AddFunc(t.Cron, func() {
		t.manager.Spawn(t)
	})
}

func concatStrings(strs ...string) string {
	var buffer bytes.Buffer

	for _, str := range strs {
		buffer.WriteString(str)
	}
	return buffer.String()
}

func (t *ProcessTemplate) NewProcessWithTrigger(trig *Trigger) *Process {
	// Replace an arg that's :event with the json representation of the event data

	args := make([]string, 0)

	for _, arg := range t.Args {
		if arg == ":json" {
			marshaled, _ := json.Marshal(trig.Data)
			args = append(args, string(marshaled))
		} else if arg == ":flags" {
			// Only pass through flags if they're strings
			for k, v := range trig.Data {
				if str, ok := v.(string); ok {
					args = append(args, "--"+k+" "+str)
				} else {
					marshaled, _ := json.Marshal(v)

					args = append(args, "--"+k+" "+string(marshaled))
				}

			}
		} else {
			args = append(args, arg)
		}
	}

	proc := t.NewProcess()

	proc.Args = args
	return proc
}

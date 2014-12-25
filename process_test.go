package gpm

import (
	. "gopkg.in/check.v1"
	"log"
	"os"
	"time"
	// "testing"
	// "fmt"
)

func (s *TestSuite) TestRespawn(c *C) {
	temp := &ProcessTemplate{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js"},
		LogFile:      "/tmp/cronlog",
		ErrFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
	}

	proc := temp.NewProcess()

	monitor := make(chan string)

	proc.monitor = monitor
	proc.Spawn()

	count := 0
	for {
		select {
		case _ = <-monitor:
			log.Println("Got worker channel int", count)
			// Should have restarted 5 times before writing to the channel
			c.Assert(proc.Respawns, Equals, 5)
			return

		default:
			// Wait a while before we check again
			time.Sleep(100 * time.Millisecond)
		}
	}
	// fmt.Println(proc.Start())
	// fmt.Println(proc.Pid)

	// proc.Watch()
}

func (s *TestSuite) TestBasic(c *C) {
	proc := &os.ProcAttr{
		Dir: "",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
			// NewLog(p.Template.LogFile),
			// NewLog(p.Template.ErrFile),
		},
	}

	os.StartProcess("/bin/cat", []string{"sample.toml", ">", "/tmp/cronlog"}, proc)
}

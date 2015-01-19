package goproc

import (
	. "gopkg.in/check.v1"
	"time"
	// "testing"
	// "fmt"
)

func (s *TestSuite) TestRespawn(c *C) {
	temp := &ProcessTemplate{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js"},
		LogFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
		Name:         "Test",
	}

	proc := temp.NewProcess()

	monitor := make(chan *Process)

	proc.monitor = monitor
	proc.Spawn()

	for {
		select {
		case _ = <-monitor:
			// Should have restarted 5 times before writing to the channel
			c.Assert(proc.Respawns, Equals, 5)
			return

		default:
			// Wait a while before we check again
			time.Sleep(100 * time.Millisecond)
		}
	}
}

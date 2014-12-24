package gpm

import (
	. "gopkg.in/check.v1"
	"time"
	// "testing"
)

func (s *TestSuite) TestSpawn(c *C) {
	proc := &Process{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js"},
		LogFile:      "/tmp/cronlog",
		ErrFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
		Name:         "test",
	}

	manager := NewManager(&Config{})
	// manager.monitor = make(chan string)
	manager.Spawn(proc)

	c.Assert(manager.Processes["test"], Equals, proc)
	// manager.Status()
}

// Make sure that a cron runs and the monitor channel receives the "finish" message
func (s *TestSuite) TestRegisterCrons(c *C) {
	proc1 := &Process{
		Command: "/usr/local/bin/node",
		Args: []string{
			"samples/longrunning.js",
		},
		Cron: "* * * * * *",
		Name: "asdf",
	}

	config := &Config{
		Cron: []*Process{
			proc1,
		},
	}

	manager := NewManager(config)

	manager.Start()

	// Make the crons are registered
	c.Assert(len(manager.cron.Entries()), Equals, 1)

	for {
		select {
		case proc := <-manager.monitor:
			// fmt.Println("Got ", proc)
			c.Assert(proc, Equals, "asdf")
			return
		default:
			// Wait a while before we check again
			time.Sleep(100 * time.Millisecond)
		}
	}

}

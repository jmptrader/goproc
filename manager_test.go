package goproc

import (
	. "gopkg.in/check.v1"
	"time"
	// "testing"
)

func (s *TestSuite) TestSpawn(c *C) {
	temp := &ProcessTemplate{
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
	go manager.Spawn(temp)

	for {
		select {
		case proc := <-manager.monitor:
			// fmt.Println("Got ", proc)
			c.Assert(proc.Template, Equals, temp)
			return
		default:
			// Wait a while before we check again
			time.Sleep(100 * time.Millisecond)
		}
	}

	// c.Assert(len(manager.Running), Equals, 1)
	// manager.Status()
}

func (s *TestSuite) TestSpawnWithLimit(c *C) {
	temp := &ProcessTemplate{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js"},
		LogFile:      "/tmp/cronlog",
		ErrFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
		Name:         "test",
	}

	manager := NewManager(&Config{
		MaxConcurrent: 1,
	})
	// manager.monitor = make(chan string)
	manager.Spawn(temp)
	manager.Spawn(temp)

	c.Assert(len(manager.Running), Equals, 1)
	c.Assert(len(manager.Queue), Equals, 1)

	go manager.Start()

	// Trigger "done"
	manager.monitor <- manager.Running[0]

	// Needs to be more than the wait time in manager.Start()
	time.Sleep(200 * time.Millisecond)
	c.Assert(len(manager.Running), Equals, 1)
	c.Assert(len(manager.Queue), Equals, 0)
	// manager.Status()
}

// Make sure that a cron runs and the monitor channel receives the "finish" message
func (s *TestSuite) TestRegisterCrons(c *C) {
	temp1 := &ProcessTemplate{
		Command: "/usr/local/bin/node",
		Args: []string{
			"samples/longrunning.js",
		},
		Cron: "* * * * * *",
		Name: "asdf",
	}

	config := &Config{
		Cron: []*ProcessTemplate{
			temp1,
		},
	}

	manager := NewManager(config)

	// log.Println(manager.monitor)
	go manager.Start()

	// Make the crons are registered

	for {
		select {
		case proc := <-manager.monitor:
			// fmt.Println("Got ", proc)
			c.Assert(proc.Template, Equals, temp1)
			return
		default:
			// Wait a while before we check again
			time.Sleep(100 * time.Millisecond)
		}
	}

}

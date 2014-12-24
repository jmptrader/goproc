package gpm

import (
	. "gopkg.in/check.v1"
	"time"
	// "testing"
)

func (s *TestSuite) TestParse(c *C) {
	procs, err := parseConfigFile("sample.toml")

	c.Assert(err, Equals, nil)

	c.Assert(len(procs.Cron), Equals, 1)
	c.Assert(len(procs.Boot), Equals, 1)
	c.Assert(len(procs.Event), Equals, 1)

	cron := procs.Cron[0]

	c.Assert(cron.Command, Equals, "php")
	c.Assert(cron.Name, Equals, "CronProcess")
	c.Assert(cron.Cron, Equals, "*")

	// proc2 := procs[1]

	// c.Assert(proc2.Command, Equals, "node")
	// // c.Assert(proc2.AutoStart, Equals, true)
	// c.Assert(proc2.KeepAlive, Equals, true)
}

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

	c.Assert(manager.processes["test"], Equals, proc)
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

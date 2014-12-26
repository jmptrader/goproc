package goproc

import (
	. "gopkg.in/check.v1"
	// "time"
	// "testing"
)

func (s *TestSuite) TestParse(c *C) {
	procs, err := ParseFile("sample.toml")

	c.Assert(err, Equals, nil)

	c.Assert(len(procs.Process), Equals, 3)

	cron := procs.Process[0]

	c.Assert(cron.Command, Equals, "php")
	c.Assert(cron.Name, Equals, "CronProcess")
	c.Assert(cron.Cron, Equals, "*")

	// proc2 := procs[1]

	// c.Assert(proc2.Command, Equals, "node")
	// // c.Assert(proc2.AutoStart, Equals, true)
	// c.Assert(proc2.KeepAlive, Equals, true)
}

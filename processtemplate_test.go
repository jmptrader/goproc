package goproc

import (
	. "gopkg.in/check.v1"
	// "log"
	// "testing"
)

func (s *TestSuite) TestNewProcessWithEventJson(c *C) {
	template := &ProcessTemplate{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js", ":json"},
		LogFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
	}

	trig := &Trigger{
		Name: "Foo",
		Data: map[string]interface{}{
			"foo":  "bar",
			"baz":  1,
			"bing": true,
		},
	}

	proc := template.NewProcessWithTrigger(trig)

	c.Assert(proc.Args[1], Equals, "{\"baz\":1,\"bing\":true,\"foo\":\"bar\"}")

}

func (s *TestSuite) TestNewProcessWithEventFlags(c *C) {
	template := &ProcessTemplate{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js", ":flags"},
		LogFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
	}

	trig := &Trigger{
		Name: "Foo",
		Data: map[string]interface{}{
			"foo": "bar",
			"baz": "bing",
		},
	}

	proc := template.NewProcessWithTrigger(trig)
	// This switches sometimes, so we allow for both orders
	if proc.Args[1] == "--foo bar" {
		c.Assert(proc.Args[2], Equals, "--baz bing")
	} else {
		c.Assert(proc.Args[1], Equals, "--baz bing")
		c.Assert(proc.Args[2], Equals, "--foo bar")

	}

}

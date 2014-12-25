package gpm

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
		ErrFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
	}

	event := &Event{
		Name: "Foo",
		Data: &map[string]interface{}{
			"foo":  "bar",
			"baz":  1,
			"bing": true,
		},
	}

	proc := template.NewProcessWithEvent(event)

	c.Assert(proc.Args[1], Equals, "{\"baz\":1,\"bing\":true,\"foo\":\"bar\"}")

}

func (s *TestSuite) TestNewProcessWithEventFlags(c *C) {
	template := &ProcessTemplate{
		Command:      "/usr/local/bin/node",
		Args:         []string{"samples/longrunning.js", ":flags"},
		LogFile:      "/tmp/cronlog",
		ErrFile:      "/tmp/cronlog",
		KeepAlive:    true,
		RespawnLimit: 5,
	}

	event := &Event{
		Name: "Foo",
		Data: &map[string]interface{}{
			"foo": "bar",
			"baz": "bing",
		},
	}

	proc := template.NewProcessWithEvent(event)
	c.Assert(proc.Args[1], Equals, "--foo \"bar\"")
	c.Assert(proc.Args[2], Equals, "--baz \"bing\"")

}

Process manager for go.

# What does it do?
Configure processes to run on a cron, at boot time, or responding to an event.

# How to use
```go

// Parse config from a toml file (see sample.toml)
config, err := gpm.ParseFile("/path/to/file.toml")
if err != nil {
	panic(err)
}
manager := gpm.NewManager(config)

// Start the processes
manager.Start()

// Access processes by name
myProc := manager.Processes["MyProcName"]
```

### Trigger an Event
```go
manager.TriggerEvent(&gpm.Event{
	Name:"My Event",
	Data:&map[string]interface{
		"foo":"bar",
	},
})
```

Data does nothing now but will eventually be accessible via special entries in the process args array, like ":data.foo" or whatever. Coming soon.


# Process Configuration
```go
type Process struct {
	// Array of arguments to pass to the command
	Args         []string

	// Executable to run (e.g. "/usr/bin/php")
	Command      string

	// Cron schedule
	Cron         string

	// Working directory (defaults to current)
	Cwd          string

	// File for process stderr
	ErrFile      string

	// Name of event to respond to
	Event        string

	// Restart this process if it dies
	KeepAlive    bool

	// File for process stdout
	LogFile      string

	// Process name (for external referencing)
	Name         string

	// Max number of respawns before stopping
	RespawnLimit int

	// # of times the process has been respawned
	Respawns     int

	// Time of the last run
	StartTime    time.Time

	// Process status ("restarting", "stopped", "restarted","exited","finished","killed","started","running")
	Status       string
}
```



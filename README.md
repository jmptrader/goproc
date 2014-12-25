# What does it do?
Configure processes to run on a cron, at boot time, or responding to an event. This is not a standalone executable - rather, you can use it within the context of your own application or build your own interaction layer.

# How to use
```go

// Parse config from a toml file (see sample.toml)
config, err := gpm.ParseFile("/path/to/file.toml")
if err != nil {
	panic(err)
}

// Or create an instance of gpm.Config from scratch (see below)
config := &gpm.Config{}


manager := gpm.NewManager(config)

// Start the processes. Note that this needs to be run in a go routine since the manager listens for finished processes to manage its queue.
go manager.Start()

// Access processes by name
myProc := manager.Processes["MyProcName"]
```

### Available Configuration Options
See "Process Template Configuration" section below for how to make a `ProcessTemplate`.
```go
type Config struct {
	// Processes to register to cron jobs
	Cron          []*ProcessTemplate

	// Processes to run on an event
	Event         []*ProcessTemplate

	// Processes to run as soon as `manager.Start()` is run
	Boot          []*ProcessTemplate

	// Most concurrent processes to run at one time. 
	// If this limit is reached, `manager.Spawn()` and `manager.TriggerEvent()` 
	// will place processes in the manager's `Queue`, which gets emptied 
	// as processes finish executing.
	MaxConcurrent int
}
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

### Special Event Arguments
You can pass event data to your process as one argument in JSON format, or as a list of key-value pairs using flags. For example, `Args:[]string{"firstArg, ":json", ":flags"}`, with the above event, would call the process with the following arguments:

```go
["firstArg", "{\"foo\":\"bar\"}", "--foo \"bar\""]
```

Note that the `:flags` indicator uses the map key as the flag name, and `json.Marshal` to determine the value.


# Process Template Configuration
Processes are run using instances of `ProcessTemplate`. 
```go
type ProcessTemplate struct {
	// Array of arguments to pass to the command.
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
}
```





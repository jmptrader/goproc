# What does it do?
Configure processes to run on a cron, at boot time, or triggered manually. This is not a standalone executable - rather, you can use it within the context of your own application or build your own interaction layer.

# How to use
```go

// Parse config from a toml file (see sample.toml)
config, err := goproc.ParseFile("/path/to/file.toml")
if err != nil {
	panic(err)
}

// Or create an instance of goproc.Config from scratch (see below)
config := &goproc.Config{}


manager := goproc.NewManager(config)

// Start the processes. Note that this needs to be run in a go routine since the manager listens for finished processes to manage its queue.
go manager.Start()

// Access processes by name
myProc := manager.Processes["MyProcName"]
```

### Available Configuration Options
See "Process Template Configuration" section below for how to make a `ProcessTemplate`.
```go
type Config struct {
	// List of process templates
	Process          []*ProcessTemplate

	// Most concurrent processes to run at one time. 
	// If this limit is reached, `manager.Spawn()` and `manager.Trigger()` 
	// will place processes in the manager's `Queue`, which gets emptied 
	// as processes finish executing.
	MaxConcurrent int
}
```

### Manually Trigger a Process
Manually trigger a process specifying name and data:

```go
manager.Trigger(&goproc.Trigger{
	Name:"My Process",
	Data:&map[string]interface{
		"foo":"bar",
	},
})
```

### Special Manual Process Arguments
You can pass custom data to your process as one argument in JSON format, or as a list of key-value pairs using flags. For example, `Args:[]string{"firstArg, ":json", ":flags"}`, with the above trigger, would call the process named "My Process" with the following arguments:

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

	// Whether to start this process as soon as the process manager starts
	AuthStart    bool

	// Executable to run (e.g. "/usr/bin/php")
	Command      string

	// Cron schedule
	Cron         string

	// Working directory (defaults to current)
	Cwd          string

	// File for process stderr
	ErrFile      string

	// Automatically restart this process if it dies
	KeepAlive    bool

	// File for process stdout
	LogFile      string

	// Process name (for external referencing)
	Name         string

	// Max number of respawns (from KeepAlive) before stopping
	RespawnLimit int
}
```





package gpm

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/robfig/cron"
	// "log"
)

type Config struct {
	Cron  []*Process
	Event []*Process
	Boot  []*Process
}

type Event struct {
	Name string
	Data map[string]interface{}
}

type Status struct {
	Processes map[string]*Process
}

func parseConfigFile(file string) (*Config, error) {

	procs := &Config{}

	_, err := toml.DecodeFile(file, procs)

	return procs, err
}

type Manager struct {
	Config    *Config
	cron      *cron.Cron
	monitor   chan string
	event     chan *Event
	processes map[string]*Process
}

func (m *Manager) Status() {
	marshaled, _ := json.Marshal(m.processes)
	fmt.Println(string(marshaled))
}

func (m *Manager) Spawn(p *Process) error {
	m.processes[p.Name] = p
	p.monitor = m.monitor
	p.Spawn()
	return nil
}

func NewManager(config *Config) *Manager {
	manager := &Manager{}
	manager.monitor = make(chan string)
	manager.processes = make(map[string]*Process)
	manager.event = make(chan *Event)
	manager.Config = config
	return manager
}

func (m *Manager) Start() {
	// Boot all processes that are set to boot
	for _, proc := range m.Config.Boot {
		m.Spawn(proc)
	}

	// Register crons
	c := cron.New()
	for _, proc := range m.Config.Cron {
		c.AddFunc(proc.Cron, func() {
			m.Spawn(proc)
		})
	}
	m.cron = c
	m.StartCrons()
}

func (m *Manager) StartCrons() {
	m.cron.Start()
}

func (m *Manager) StopCrons() {
	m.cron.Stop()
}

func (m *Manager) TriggerEvent(event *Event) {
	// Loop through all event processes to see which ones respond
	for _, proc := range m.Config.Event {
		if proc.Event == event.Name {
			m.Spawn(proc)
		}
	}
}

// func ()

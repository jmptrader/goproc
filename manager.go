package gpm

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
)

type Manager struct {
	Config    *Config
	cron      *cron.Cron
	monitor   chan string
	event     chan *Event
	Processes map[string]*ProcessTemplate
}

func (m *Manager) Status() {
	marshaled, _ := json.Marshal(m.Processes)
	fmt.Println(string(marshaled))
}

func (m *Manager) Spawn(t *ProcessTemplate) error {
	// Create a new process from template
	p := t.NewProcess()

	m.Processes[t.Name] = t
	p.monitor = m.monitor
	p.Spawn()
	return nil
}

func NewManager(config *Config) *Manager {
	manager := &Manager{}
	manager.monitor = make(chan string)
	manager.Processes = make(map[string]*ProcessTemplate)
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

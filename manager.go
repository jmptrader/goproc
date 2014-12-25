package goproc

import (
	// "encoding/json"
	// "fmt"
	"github.com/robfig/cron"
	"time"
)

type Manager struct {
	Config  *Config
	cron    *cron.Cron
	monitor chan *Process
	event   chan *Event
	Queue   []*Process
	Running []*Process
}

func (m *Manager) Status() {
	// marshaled, _ := json.Marshal(m.Processes)
	// fmt.Println(string(marshaled))
}

func (m *Manager) Spawn(t *ProcessTemplate) {
	// Create a new process from template
	p := t.NewProcess()
	m.spawn(p)
}

func (m *Manager) TriggerEvent(event *Event) {
	// Loop through all event processes to see which ones respond
	for _, t := range m.Config.Event {
		if t.Event == event.Name {
			p := t.NewProcessWithEvent(event)
			m.spawn(p)
		}
	}
}

func (m *Manager) spawn(p *Process) {
	p.monitor = m.monitor

	// Are we at the limit?
	if m.Config.MaxConcurrent > 0 && len(m.Running) >= m.Config.MaxConcurrent {
		m.Queue = append(m.Queue, p)
	} else {
		m.Running = append(m.Running, p)
		p.Spawn()
	}
}

func NewManager(config *Config) *Manager {
	manager := &Manager{}
	manager.monitor = make(chan *Process)
	manager.Queue = make([]*Process, 0)
	manager.Running = make([]*Process, 0)
	manager.event = make(chan *Event)
	manager.Config = config
	return manager
}

func (m *Manager) Start() {
	// Boot all processes that are set to boot
	for _, t := range m.Config.Boot {
		m.Spawn(t)
	}

	// Register crons
	c := cron.New()
	for _, t := range m.Config.Cron {
		c.AddFunc(t.Cron, func() {
			m.Spawn(t)
		})
	}
	m.cron = c
	m.StartCrons()

	for {
		select {
		case proc := <-m.monitor:
			// Remove proc from running
			for i, p := range m.Running {
				if proc == p {
					m.Running = append(m.Running[:i], m.Running[i+1:]...)

					// Pop one off the queue
					if len(m.Queue) > 0 {
						next := m.Queue[0]

						if len(m.Queue) == 1 {
							m.Queue = make([]*Process, 0)
						} else {
							m.Queue = append([]*Process{}, m.Queue[1:]...)
						}
						m.spawn(next)
					}
				}
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (m *Manager) StartCrons() {
	m.cron.Start()
}

func (m *Manager) StopCrons() {
	m.cron.Stop()
}

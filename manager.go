package goproc

import (
	// "encoding/json"
	"github.com/robfig/cron"
	"log"
	"time"
)

type Manager struct {
	Config  *Config
	cron    *cron.Cron
	monitor chan *Process
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

func (m *Manager) Trigger(trig *Trigger) {
	// Loop through all event processes to see which ones respond
	for _, t := range m.Config.Process {
		if t.Name == trig.Name {
			p := t.NewProcessWithTrigger(trig)
			m.spawn(p)
		}
	}
}

func (m *Manager) spawn(p *Process) {
	p.monitor = m.monitor
	// // Are we at the limit?
	if m.Config.MaxConcurrent > 0 && len(m.Running) >= m.Config.MaxConcurrent {
		p.QueuedAt = time.Now()
		log.Println("Queuing process ", p.Template.Name)
		m.Queue = append(m.Queue, p)
	} else {
		log.Println("Booting process ", p.Template.Name)
		m.Running = append(m.Running, p)
		p.Spawn()
	}
}

func NewManager(config *Config) *Manager {
	manager := &Manager{}
	manager.monitor = make(chan *Process)
	manager.Queue = make([]*Process, 0)
	manager.Running = make([]*Process, 0)
	manager.Config = config
	return manager
}

func (m *Manager) Start() {
	log.Println("Starting process manager")
	c := cron.New()

	for _, t := range m.Config.Process {
		if t.AutoStart {
			log.Printf("Booting %s\n", t.Name)
			m.Spawn(t)
		}

		if len(t.Cron) > 0 {
			log.Printf("Adding %s to crontab\n", t.Name)
			c.AddFunc(t.Cron, func() {
				m.Spawn(t)
			})
		}
	}

	m.cron = c
	m.StartCrons()

	for {
		select {
		case proc := <-m.monitor:
			log.Println("Got proc from monitor channel", proc.Template.Name)
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

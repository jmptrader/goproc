package goproc

import (
	// "encoding/json"
	"errors"
	"github.com/robfig/cron"
	"github.com/twinj/uuid"
	"log"
	"time"
)

type Manager struct {
	Config  *Config
	cron    *cron.Cron
	monitor chan *Process
	Queue   []*Process
	Running []*Process
	Logger  Logger
	Started bool
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type BasicLogger struct{}

func (b *BasicLogger) Info(msg string) {
	log.Println(msg)
}

func (b *BasicLogger) Error(msg string) {
	log.Println(msg)
}

func (m *Manager) Status() {
	// marshaled, _ := json.Marshal(m.Processes)
	// fmt.Println(string(marshaled))
}

func (m *Manager) Spawn(t *ProcessTemplate) {
	// Create a new process from template
	p := t.NewProcess()
	m.SpawnProcess(p)
}

func (m *Manager) Trigger(trig *Trigger) (*Process, error) {
	// Loop through all event processes to see which ones respond
	p, err := m.ProcessFromTrigger(trig)
	if err == nil {
		m.Logger.Info("Found process. Spawning")
		m.SpawnProcess(p)
	} else {
		m.Logger.Error("No matching process found: " + trig.Name)
	}
	return p, err
}

func (m *Manager) ProcessFromTrigger(trig *Trigger) (*Process, error) {
	for _, t := range m.Config.Process {
		if t.Name == trig.Name {
			p := t.NewProcessWithTrigger(trig)
			return p, nil
		}
	}
	return &Process{}, errors.New("No process responds to this trigger")
}

func (m *Manager) SpawnProcess(p *Process) {
	p.monitor = m.monitor

	// // Are we at the limit?
	if m.Config.MaxConcurrent > 0 && len(m.Running) >= m.Config.MaxConcurrent {
		p.QueuedAt = time.Now()
		m.Logger.Info("Queuing process " + p.Template.Name)
		p.Status = "queued"
		m.Queue = append(m.Queue, p)

		for _, hook := range m.Config.ProcessHooks {
			hook("queued", p)
		}
	} else {
		m.Logger.Info("Booting process " + p.Template.Name)
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
	manager.Logger = &BasicLogger{}

	// log.Println(config.ProcessHooks)
	return manager
}

func (m *Manager) Start() {
	uuid.SwitchFormat(uuid.CleanHyphen)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	m.Logger.Info("Starting process manager")

	c := cron.New()

	for _, t := range m.Config.Process {
		t.manager = m
		if t.AutoStart {
			m.Logger.Info("Booting " + t.Name)
			m.Spawn(t)
		}

		if len(t.Cron) > 0 {
			m.Logger.Info("Adding " + t.Name + " to crontab")
			t.RegisterCron(c)
		}
	}

	m.cron = c
	m.Logger.Info("Created cron...")
	// m.StartCrons()
	m.Started = true

	m.Logger.Info("Set started to true")
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
						m.SpawnProcess(next)
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

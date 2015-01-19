package goproc

import (
	// "labix.org/v2/mgo/bson"
	"fmt"
	// "io"

	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"time"
)

type Process struct {
	Template  *ProcessTemplate
	x         *os.Process
	Args      []string
	Hostname  string
	Pid       int
	Status    string
	monitor   chan *Process
	StartTime time.Time
	Respawns  int
	Info      *os.ProcessState
	Error     error
	QueuedAt  time.Time
	Uuid      string
	LogFile   string
	// LogWriter io.Writer
}

func NewLog(path string, reset bool) *os.File {
	if path == "" {
		return nil
	}

	flag := os.O_CREATE | os.O_RDWR | os.O_APPEND
	if reset {
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND | os.O_TRUNC
	}

	folder := filepath.Dir(path)

	_, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(folder, 0775)
		}
	}

	file, err := os.OpenFile(path, flag, 0660)
	if err != nil {
		log.Println("Failed to create log file!!", err.Error())
		return nil
	}
	return file
}
func (p *Process) release(status string) {
	// Hooks?
	if p.hasManager() {
		for _, hook := range p.Template.manager.Config.ProcessHooks {
			hook(status, p)
		}
	}

	if p.x != nil {
		p.x.Release()
	}
	p.Pid = 0
	// p.Pidfile.delete()
	p.Status = status

	if status != "restarting" {

		go func() {
			p.monitor <- p
		}()
	}

}

func (p *Process) Watch() {
	if p.x == nil {
		p.release("stopped")
		return
	}
	status := make(chan *os.ProcessState)
	died := make(chan error)
	go func() {
		state, err := p.x.Wait()
		if err != nil {
			died <- err
			return
		}
		status <- state
	}()
	select {
	case state := <-status:
		if p.Status == "stopped" {
			return
		}
		p.Info = state

		if p.Template.KeepAlive {

			if p.hasManager() {
				p.Template.manager.Logger.Info("Process " + p.Template.Name + " died unexpectedly - restarting")
			}
			if p.Template.RespawnLimit > 0 && p.Respawns == p.Template.RespawnLimit {
				p.Stop()
			} else {
				p.Respawns++
				p.Status = "restarted"
				p.Restart()

			}

			return
		} else {
			if state.Success() {
				p.release("finished")
			} else {
				p.release("error")
			}

		}

	case err := <-died:
		p.Error = err
		p.release("error")
	}
}

func (p *Process) Start() bool {
	p.Hostname, _ = os.Hostname()
	p.Template.RunCount++
	p.StartTime = time.Now()
	// wd, _ := os.Getwd()

	// Filter log file?
	fileName := p.Template.LogFile

	if p.hasManager() {
		for _, filter := range p.Template.manager.Config.LogFilters {
			fileName = filter(fileName, p)
		}
	}

	logFile := NewLog(fileName, p.Template.ResetLog)
	p.LogFile = fileName
	// Hooks?
	if p.hasManager() {
		for _, hook := range p.Template.manager.Config.ProcessHooks {
			hook("start", p)
		}
	}

	proc := &os.ProcAttr{
		Dir: p.Template.Cwd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			// os.Stdout,
			// os.Stderr,
			logFile,
			logFile,
		},
	}
	args := append([]string{p.Template.Name}, p.Args...)

	// process, err := os.StartProcess("/bin/cat", []string{"sample.toml"}, proc)
	process, err := os.StartProcess(p.Template.Command, args, proc)
	if err != nil {
		if p.hasManager() {
			p.Template.manager.Logger.Error(fmt.Sprintf("%s failed. %s\n", p.Template.Name, err.Error()))
		}
		return false
	}
	// err = p.Pidfile.write(process.Pid)
	// if err != nil {
	// 	log.Printf("%s pidfile error: %s\n", p.Name, err)
	// 	return ""
	// }

	p.x = process
	p.Pid = process.Pid
	p.Status = "started"
	return true
}

func (p *Process) Restart() {
	// Hooks?
	if p.Status != "restarted" {
		p._stop()
	}

	p.release("restarting")
	p.Spawn()
}

func (p *Process) _stop() {
	if p.hasManager() {
		p.Template.manager.Logger.Info("Killing process: " + p.Template.Name)
	}
	if p.x != nil {
		// p.x.Kill() this seems to cause trouble
		cmd := exec.Command("kill", fmt.Sprintf("%d", p.x.Pid))
		o, err := cmd.CombinedOutput()

		if err != nil && p.hasManager() {
			p.Template.manager.Logger.Error(string(o))
		}
		// p.children.stop("all")
	}

}

func (p *Process) hasManager() bool {
	v := reflect.ValueOf(p.Template.manager)

	return !v.IsNil()
}

func (p *Process) Spawn() {

	// Hooks?
	if p.hasManager() {
		for _, hook := range p.Template.manager.Config.ProcessHooks {
			hook("spawn", p)
		}
	}

	go func() {
		p.Start()
		if p.Pid > 0 {
			p.Status = "running"
			for _, hook := range p.Template.manager.Config.ProcessHooks {
				hook("running", p)
			}
		}

		go p.Watch()
	}()
}

func (p *Process) Stop() string {
	p._stop()
	p.release("stopped")
	message := fmt.Sprintf("%s stopped.\n", p.Template.Name)
	return message
}

package goproc

import (
	// "labix.org/v2/mgo/bson"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	Template  *ProcessTemplate
	x         *os.Process
	Args      []string
	Pid       int
	Status    string
	monitor   chan *Process
	StartTime time.Time
	Respawns  int
	Info      *os.ProcessState
	Error     error
	QueuedAt  time.Time
}

func NewLog(path string) *os.File {
	if path == "" {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("%s", err)
		return nil
	}
	return file
}
func (p *Process) release(status string) {
	log.Printf("Releasing process %d (%s) with status %s\n", p.Pid, p.Template.Name, status)
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

			if p.Template.RespawnLimit > 0 && p.Respawns == p.Template.RespawnLimit {
				p.Stop()
				log.Printf("%s respawn limit reached.\n", p.Template.Name)
			} else {
				p.Respawns++
				p.Restart()
				p.Status = "restarted"
			}

			return
		} else {
			p.release("finished")
		}

	case err := <-died:
		p.Error = err
		p.release("killed")
		log.Printf("%d %s killed = %#v", p.x.Pid, p.Template.Name, err)
	}
}

func (p *Process) Start() bool {
	p.Template.RunCount++
	p.StartTime = time.Now()
	// wd, _ := os.Getwd()
	proc := &os.ProcAttr{
		Dir: p.Template.Cwd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			// os.Stdout,
			// os.Stderr,
			NewLog(p.Template.LogFile),
			NewLog(p.Template.ErrFile),
		},
	}
	args := append([]string{p.Template.Name}, p.Args...)
	// process, err := os.StartProcess("/bin/cat", []string{"sample.toml"}, proc)
	process, err := os.StartProcess(p.Template.Command, args, proc)
	if err != nil {
		log.Fatalf("%s failed. %s\n", p.Template.Name, err)
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
	p._stop()
	p.release("restarting")
	p.Spawn()
}

func (p *Process) _stop() {
	if p.x != nil {
		// p.x.Kill() this seems to cause trouble
		cmd := exec.Command("kill", fmt.Sprintf("%d", p.x.Pid))
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		// p.children.stop("all")
	}

}

func (p *Process) Spawn() {
	go func() {
		p.Start()
		if p.Pid > 0 {
			p.Status = "running"
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

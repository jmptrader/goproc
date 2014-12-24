package gpm

import (
	// "labix.org/v2/mgo/bson"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	Command      string
	Cwd          string
	KeepAlive    bool
	Event        string
	LogFile      string
	ErrFile      string
	Pid          int
	Status       string
	Args         []string
	Name         string
	Process      *os.Process
	Respawns     int
	RespawnLimit int
	Cron         string
	monitor      chan<- string
	StartTime    time.Time
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
	if p.Process != nil {
		p.Process.Release()
	}
	p.Pid = 0
	// p.Pidfile.delete()
	p.Status = status

	if status != "restarting" {
		p.monitor <- p.Name
	}

}

func (p *Process) Watch() {
	if p.Process == nil {
		p.release("stopped")
		return
	}
	status := make(chan *os.ProcessState)
	died := make(chan error)
	go func() {
		state, err := p.Process.Wait()
		if err != nil {
			died <- err
			return
		}
		status <- state
	}()
	select {
	case _ = <-status:
		if p.Status == "stopped" {
			return
		}

		if p.KeepAlive {

			if p.Respawns == p.RespawnLimit {
				p.release("exited")
				log.Printf("%s respawn limit reached.\n", p.Name)
			} else {
				p.Respawns++
				p.restart()
				p.Status = "restarted"
			}
			return
		} else {
			p.release("finished")
		}

	case err := <-died:
		p.release("killed")
		log.Printf("%d %s killed = %#v", p.Process.Pid, p.Name, err)
	}
}

func (p *Process) Start() bool {
	p.StartTime = time.Now()

	// wd, _ := os.Getwd()
	proc := &os.ProcAttr{
		Dir: p.Cwd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			NewLog(p.LogFile),
			NewLog(p.ErrFile),
		},
	}
	args := append([]string{p.Name}, p.Args...)
	process, err := os.StartProcess(p.Command, args, proc)
	if err != nil {
		log.Fatalf("%s failed. %s\n", p.Name, err)
		return false
	}
	// err = p.Pidfile.write(process.Pid)
	// if err != nil {
	// 	log.Printf("%s pidfile error: %s\n", p.Name, err)
	// 	return ""
	// }

	p.Process = process
	p.Pid = process.Pid
	p.Status = "started"
	return true
}

func (p *Process) restart() {
	p._stop()
	p.release("restarting")
	p.Spawn()
}

func (p *Process) _stop() {
	if p.Process != nil {
		// p.x.Kill() this seems to cause trouble
		cmd := exec.Command("kill", fmt.Sprintf("%d", p.Process.Pid))
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

func (p *Process) stop() string {
	p._stop()
	p.release("stopped")
	message := fmt.Sprintf("%s stopped.\n", p.Name)
	return message
}

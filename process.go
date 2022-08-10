package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Process struct {
	Id, Pid     int
	File, State string
	Args        []string
	Cmd         *exec.Cmd
	Logger      *log.Logger
	WorkingDir  string
}

func (p *Process) start() {
	cmd := exec.Command(p.File, p.Args...)
	cmd.Dir = p.WorkingDir
	p.Cmd = cmd
	if cmd == nil {
		fmt.Println("cmd is nil, WTF")
	}
	stdout, stderr := p.getPipes()
	err := cmd.Start()
	if err != nil {
		p.Logger.Printf("%s", p.Id, err)
		return
	}
	p.State = "running"
	go p.watchPipe(stdout)
	go p.watchPipe(stderr)

}

func (p *Process) stop() {
	if p.Cmd == nil {
		println("cmd is nil, WTF")
	}
	err := p.Cmd.Process.Signal(os.Signal(syscall.SIGTERM))
	if err != nil {
		p.Logger.Println(err.Error())
		return
	}
	p.State = "stopped"
}

func (p *Process) restart() {
	p.stop()
	p.start()
}

func (p *Process) kill() {
	p.Cmd.Process.Kill()
	p.State = "killed"
}

func (p *Process) getPipes() (io.Reader, io.Reader) {
	fmt.Println("watchProcess")
	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		p.Logger.Println(err.Error())
		return nil, nil
	}
	stderr, err := p.Cmd.StderrPipe()
	if err != nil {
		p.Logger.Println(err.Error())
		return nil, nil
	}
	return stdout, stderr

}

func (p *Process) watchPipe(pipe io.Reader) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		p.Logger.Print(scanner.Text())
	}
	p.State = "exited"
}

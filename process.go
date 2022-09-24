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
	Stdin       *io.Writer
	out         chan string
}

func (p *Process) start() {
	cmd := exec.Command(p.File, p.Args...)
	cmd.Dir = p.WorkingDir
	p.Cmd = cmd
	if cmd == nil {
		fmt.Println("cmd is nil, WTF")
	}
	stdout, stderr, stdin := p.getPipes()
	p.Stdin = &stdin
	err := cmd.Start()
	p.Pid = cmd.Process.Pid
	if err != nil {
		p.Logger.Printf("%s", p.Id, err)
		return
	}
	p.State = "running"

	p.out = make(chan string, 10)
	go p.watchPipe(stdout, p.out)
	go p.watchPipe(stderr, p.out)

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

func (p *Process) getPipes() (io.Reader, io.Reader, io.Writer) {
	fmt.Println("watchProcess")
	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		p.Logger.Println(err.Error())
		return nil, nil, nil
	}
	stderr, err := p.Cmd.StderrPipe()
	if err != nil {
		p.Logger.Println(err.Error())
		return nil, nil, nil
	}
	stdin, err := p.Cmd.StdinPipe()
	if err != nil {
		p.Logger.Println(err.Error())
		return nil, nil, nil
	}
	return stdout, stderr, stdin

}

func (p *Process) watchPipe(pipe io.Reader, c *chan string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		text := scanner.Text()
		p.Logger.Print(text)
		if len(*c) > 10 {
			*c = make(chan string, 10)
		}
		*c <- text
	}
	p.State = "exited"
}

func (p Process) stdin(input string) {
	println("writing to stdin: " + input)
	(*p.Stdin).Write([]byte(input))
}

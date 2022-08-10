package main

import (
	"log"
	"os"
	"strconv"
)

type KPMD int

func (t *KPMD) Run(args []string, reply *int) error {
	dir, _ := os.UserHomeDir()
	err := os.Mkdir(dir+"/.kpmd", 0755)
	if err != nil {
		log.Println(err.Error())
	}
	logFile, err := os.OpenFile(dir+"/.kpmd/process"+strconv.Itoa(len(processes))+".log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		repl := 1
		reply = &repl
		return err
	}
	workdir := args[len(args)-1]
	process := Process{
		Id:         len(processes),
		File:       args[1],
		Args:       args[2 : len(args)-1],
		Logger:     log.New(logFile, strconv.Itoa(len(processes))+"|\t", 0),
		WorkingDir: workdir,
	}
	processes = append(processes, &process)

	go process.start()
	return nil
}

func (t *KPMD) Stop(args []string, reply *int) error {
	id, _ := strconv.Atoi(args[0])
	for _, process := range processes {
		if process.Id == id {
			process.stop()
		}
	}
	return nil
}

func (t *KPMD) Restart(args []string, reply *int) error {
	id, _ := strconv.Atoi(args[0])
	for _, process := range processes {
		if process.Id == id {
			process.restart()
		}
	}
	return nil
}

func (t *KPMD) Kill(args []string, reply *int) error {
	id, _ := strconv.Atoi(args[0])
	for _, process := range processes {
		if process.Id == id {
			process.kill()
		}
	}
	return nil
}

func (t *KPMD) List(args string, reply *[]ProcessInfo) error {
	var procs []ProcessInfo
	for _, process := range processes {
		procs = append(procs, ProcessInfo{
			Id:    process.Id,
			Pid:   process.Pid,
			File:  process.File,
			State: process.State,
		})
	}
	*reply = procs
	return nil
}

type ProcessInfo struct {
	Id, Pid     int
	File, State string
	Args        []string
}

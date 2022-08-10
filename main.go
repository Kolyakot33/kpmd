package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var processes []*Process

func main() {

	processes = make([]*Process, 0)

	msg := new(KPMD)
	rpc.Register(msg)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":7124")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

/*
func runProcess(name string, args []string) error {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go watchProcess(stdout, cmd)
	cmd.Start()
	return nil
}

func watchProcess(stdout io.Reader, cmd *exec.Cmd) {
	for {
		buffer := make([]byte, 16)
		_, err := stdout.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Process finished")
			} else {
				log.Fatal(err)
			}
			break
		}
		print(string(buffer))
	}
}
*/

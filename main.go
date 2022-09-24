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

package main

import (
	"fmt"
	"log"
	"net"
)

const (
	network = "tcp"
)

type Port int

func (p Port) Addr() string {
	return fmt.Sprintf(":%d", p)
}

func main() {
	port := Port(9000)

	ln, err := net.Listen(network, port.Addr())
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	defer ln.Close()

	log.Println("listening on port", port)

	server := NewServer()
	go server.run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("could not accept connection: %v", err)
			continue
		}

		go server.newClient(conn)
	}
}

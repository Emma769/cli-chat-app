package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type roomer interface {
	Broadcast(clienter, string)
	Name() string
	Remove(clienter)
	Add(clienter)
}

type Server struct {
	rooms   map[string]roomer
	actions chan actioner
}

func NewServer() *Server {
	return &Server{
		rooms:   map[string]roomer{},
		actions: make(chan actioner),
	}
}

func (s *Server) newClient(conn net.Conn) {
	log.Println("creating new client from connection", conn.RemoteAddr().String())
	cl := NewClient(conn, s.actions)
	cl.readloop()
}

func (s *Server) run() {
	for action := range s.actions {
		switch action.Cmd() {
		case cmd_join:
			s.join(action.Client(), action.Args())
		case cmd_list:
			s.list(action.Client(), action.Args())
		case cmd_message:
			s.message(action.Client(), action.Args())
		case cmd_name:
			s.name(action.Client(), action.Args())
		case cmd_quit:
			s.quit(action.Client(), action.Args())
		}
	}
}

func (s *Server) quit(client clienter, _ []string) {
	roomName := client.Room().Name()
	s.rooms[roomName].Broadcast(
		client, fmt.Sprintf("%s left the room", client.Name()),
	)
	s.ExitRoom(client)
	log.Println("sad to see you go", client.Name())
	client.SetRoom(nil)
	client.Close()
}

func (s *Server) name(client clienter, args []string) {
	name := fst(args)
	client.SetName(name)
	client.Write(fmt.Sprintf("alright, I'll call you %s", name))
}

func (s *Server) message(client clienter, args []string) {
	if client.Room() == nil {
		client.Write("you must join a room")
		return
	}

	message := strings.Join(args, " ")
	roomName := client.Room().Name()
	s.rooms[roomName].Broadcast(client, message)
}

func (s *Server) list(client clienter, _ []string) {
	roomlist := []string{}

	for name := range s.rooms {
		roomlist = append(roomlist, name)
	}

	if len(roomlist) == 0 {
		client.Write("no available rooms")
		return
	}

	rooms := strings.Join(roomlist, ", ")

	client.Write(fmt.Sprintf("all available rooms: %s", rooms))
}

func (s *Server) join(client clienter, args []string) {
	roomName := fst(args)

	r, ok := s.rooms[roomName]

	if !ok {
		r = &Room{
			name:    roomName,
			members: map[string]clienter{},
		}

		s.rooms[roomName] = r
	}

	s.ExitRoom(client)
	r.Add(client)
	client.SetRoom(r)
	r.Broadcast(client, fmt.Sprintf("%s joined the room", client.Name()))
}

func (s *Server) ExitRoom(client clienter) {
	if client.Room() != nil {
		roomName := client.Room().Name()
		s.rooms[roomName].Remove(client)
		s.rooms[roomName].Broadcast(
			client, fmt.Sprintf("%s left the room", client.Name()),
		)
	}
}

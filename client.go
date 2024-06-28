package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type actioner interface {
	Cmd() cmd
	Args() []string
	Client() clienter
}

type Client struct {
	id     string
	name   string
	conn   net.Conn
	room   roomer
	action chan<- actioner
}

func NewClient(conn net.Conn, action chan<- actioner) *Client {
	return &Client{
		id:     GenID(),
		name:   "anonymous",
		conn:   conn,
		action: action,
	}
}

func (c Client) Room() roomer {
	return c.room
}

func (c *Client) SetRoom(room roomer) {
	c.room = room
}

func (c Client) ID() string {
	return c.id
}

func (c *Client) SetName(name string) {
	c.name = name
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) WriteErr(err error) {
	c.conn.Write([]byte(fmt.Sprintf("ERR: %s \n", err.Error())))
}

func (c *Client) Write(msg string) {
	c.conn.Write([]byte(fmt.Sprintf("> %s\n", msg)))
}

func (c *Client) readloop() {
	for {
		reader := bufio.NewReader(c.conn)

		input, err := reader.ReadString('\n')

		if err != nil && !errors.Is(err, io.EOF) {
			log.Fatal(err)
		}

		input = strings.Trim(input, "\r\n")

		f, rest := fstnrst(strings.Fields(input))

		switch f {
		case "/name":
			c.action <- Action{
				cmd:    cmd_name,
				client: c,
				args:   rest,
			}
		case "/join":
			c.action <- Action{
				cmd:    cmd_join,
				client: c,
				args:   rest,
			}
		case "/list":
			c.action <- Action{
				cmd:    cmd_list,
				client: c,
				args:   rest,
			}
		case "/quit":
			c.action <- Action{
				cmd:    cmd_quit,
				client: c,
				args:   rest,
			}
		case "/msg":
			c.action <- Action{
				cmd:    cmd_message,
				client: c,
				args:   rest,
			}
		default:
			c.WriteErr(fmt.Errorf("unknown command: %s", f))
		}
	}
}

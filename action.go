package main

type cmd int

const (
	cmd_name cmd = iota
	cmd_quit
	cmd_join
	cmd_list
	cmd_message
)

type Action struct {
	cmd    cmd
	args   []string
	client clienter
}

func (a Action) Cmd() cmd {
	return a.cmd
}

func (a Action) Args() []string {
	return a.args
}

func (a Action) Client() clienter {
	return a.client
}

package main

type clienter interface {
	ID() string
	SetName(string)
	SetRoom(roomer)
	Room() roomer
	Write(string)
	Name() string
	Close() error
}

type Room struct {
	name    string
	members map[string]clienter
}

func (r Room) Name() string {
	return r.name
}

func (r *Room) Broadcast(client clienter, message string) {
	for id, member := range r.members {
		if id != client.ID() {
			member.Write(message)
		}
	}
}

func (r *Room) Remove(client clienter) {
	delete(r.members, client.ID())
}

func (r *Room) Add(client clienter) {
	r.members[client.ID()] = client
}

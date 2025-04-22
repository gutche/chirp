package main

import "time"

type Room struct {
	ID          string
	Clients     map[*Client]bool
	Broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	TimeoutFunc func()
	CreatedAt time.Time
	EndsAt int64
}

func NewRoom(id string, endsAt int64, onExpire func()) *Room {
	room := &Room{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		TimeoutFunc: onExpire,
		EndsAt:     endsAt,
		CreatedAt:  time.Now(),
	}

	duration := time.Until(time.Unix(endsAt, 0))

	if duration <= 0 {
		go onExpire()
		return room
	}

	time.AfterFunc(duration, func() {
		onExpire()
		room.CloseAll()
	})

	return room
}

func (r *Room) CloseAll() {
	for client := range r.Clients {
		close(client.Send)
		client.Conn.Close()
		delete(r.Clients, client)
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

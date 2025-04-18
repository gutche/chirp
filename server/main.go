package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var rooms = make(map[string]*Room) // room ID -> Room

func main() {
	http.HandleFunc("/ws", handleConnections)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for simplicity
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	username := r.URL.Query().Get("username")

	if roomID == "" || username == "" {
		http.Error(w, "Missing room or username", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	room, ok := rooms[roomID]
	if !ok {
		room = NewRoom(roomID)
		rooms[roomID] = room
		go room.Run() // Start the room
	}

	client := &Client{
		Username: username,
		Conn:     conn,
		Room:     room,
		Send:     make(chan []byte),
	}

	room.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

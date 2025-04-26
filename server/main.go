package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var rooms = make(map[string]*Room) // room ID -> Room

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/ws/check", withCORS(checkRoom))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for simplicity
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomID")
	username := r.URL.Query().Get("username")
	expiration := r.URL.Query().Get("expiration")
	create := r.URL.Query().Get("create") == "true"
	log.Println(r.URL.Query())

	if roomID == "" || username == "" {
		http.Error(w, "Missing room or username", http.StatusBadRequest)
		return
	}
	room, ok := rooms[roomID]
	if !ok {
		if !create {
			http.Error(w, "Room does not exist", http.StatusNotFound)
			return
		}

		if expiration == "" {
			http.Error(w, "Missing expiration timestamp", http.StatusBadRequest)
			return
		}

		endsAt, err := strconv.ParseInt(expiration, 10, 64)
		if err != nil {
			http.Error(w, "Invalid expiration timestamp", http.StatusBadRequest)
			return
		}

		room = NewRoom(roomID, endsAt, func() {
			delete(rooms, roomID)
			log.Printf("Room %s expired and deleted\n", roomID)
		})
		rooms[roomID] = room
		go room.Run()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		Username: username,
		Conn:     conn,
		Room:     room,
		Send:     make(chan []byte),
	}

	initMsg := map[string]interface{}{
		"type":   "init",
		"endsAt": room.EndsAt,
	}

	if err := conn.WriteJSON(initMsg); err != nil {
		log.Println("Error sending init message:", err)
		conn.Close()
		return
	}

	room.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

func checkRoom(w http.ResponseWriter, r *http.Request) {
    roomID := r.URL.Query().Get("roomID")
    if roomID == "" {
        http.Error(w, "Missing roomID", http.StatusBadRequest)
        return
    }
    if _, ok := rooms[roomID]; !ok {
        http.Error(w, "Room does not exist", http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        h(w, r)
    }
}

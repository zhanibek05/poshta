package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID int64
	ChatID int64
	Send   chan []byte
}

type Hub struct {
	clients    map[int64]map[*Client]bool // chatID -> clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastPayload
	mu         sync.Mutex
}

type BroadcastPayload struct {
	ChatID int64
	Data   []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastPayload),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.ChatID] == nil {
				h.clients[client.ChatID] = make(map[*Client]bool)
			}
			h.clients[client.ChatID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.ChatID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
				}
			}
			h.mu.Unlock()

		case payload := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.clients[payload.ChatID]; ok {
				for client := range clients {
					select {
					case client.Send <- payload.Data:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

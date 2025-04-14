package websocket


import (
	"encoding/json"
	"net/http"
	"poshta/pkg/reqresp"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDStr := mux.Vars(r)["chat_id"]
		userIDStr := r.URL.Query().Get("user_id")
		chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &Client{
			Conn:   conn,
			UserID: userID,
			ChatID: chatID,
			Send:   make(chan []byte),
		}

		hub.register <- client

		go client.readPump(hub)
		go client.writePump()
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Optionally decode and store in DB
		var msg reqresp.SendMessageRequest
		if err := json.Unmarshal(msgBytes, &msg); err == nil {
			msg.SenderID = c.UserID
			msg.ChatID = c.ChatID
			msgBytes, _ = json.Marshal(msg)
		}

		hub.broadcast <- BroadcastPayload{
			ChatID: c.ChatID,
			Data:   msgBytes,
		}
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

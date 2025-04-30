package handlers

import (
	"net/http"
	"poshta/internal/app/ws"
	"poshta/internal/usecase"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	Hub            *ws.Hub
	MessageUseCase usecase.MessageUseCase
	ChatUseCase    usecase.ChatService
}

func NewWSHandler(hub *ws.Hub, msgUC usecase.MessageUseCase, chatUC usecase.ChatService) *WSHandler {
	return &WSHandler{
		Hub:            hub,
		MessageUseCase: msgUC,
		ChatUseCase:    chatUC,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade", http.StatusInternalServerError)
		return
	}

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
		Hub:    h.Hub,
		Send:   make(chan []byte, 256),
	}

	h.Hub.Register <- client
	go client.WritePump()
	go client.ReadPump(h.MessageUseCase, h.ChatUseCase)
}

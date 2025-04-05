package handlers

import (
	"poshta/internal/service"
)

type ChatHandler struct {
	chatService service.Chatservice
}

func NewChatHandler(chatService service.Chatservice) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

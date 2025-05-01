package ws

import (
	"context"
	"encoding/json"
	"poshta/internal/domain/models"
	"poshta/internal/usecase"
	"poshta/pkg/reqresp"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan []byte
}

func (c *Client) ReadPump(messageUseCase usecase.MessageUseCase, chatUseCase usecase.ChatService) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg models.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		// сохраняем в БД
		_, _ = messageUseCase.SendMessage(context.Background(), reqresp.SendMessageRequest{
			ChatID:   msg.ChatID,
			SenderID: msg.SenderID,
			Content:  msg.Content,
			EncryptedKey: msg.EncryptedKey,
		})



		// получаем чат и второго пользователя
		chat, _ := chatUseCase.GetChatByID(context.Background(), msg.ChatID)
		recipient := chat.User1ID
		if recipient == msg.SenderID {
			recipient = chat.User2ID
		}

		c.Hub.SendTo <- TargetedMessage{
			RecipientIDs: []string{msg.SenderID, recipient},
			Message:      msgBytes,
		}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

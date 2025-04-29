package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"poshta/internal/domain/models"
	"poshta/internal/middleware"
	"poshta/internal/usecase"
	"poshta/pkg/reqresp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	messageUseCase usecase.MessageUseCase
}

func NewMessageHandler(messageUseCase usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
	}
}

// @Summary Send a message
// @Description Send a message to a chat
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body reqresp.SendMessageRequest true "Create message request"
// @Success 201 {object} models.Message "Chat created successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid request"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /message [post]
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req reqresp.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	message, err := h.messageUseCase.SendMessage(r.Context(), req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, message)
}



// DeleteMessage godoc
// @Summary      Delete a message
// @Description  Deletes a message by ID. Only the owner of the message can delete it.
// @Tags         messages
// @Security     BearerAuth
// @Param        id   path      int  true  "Message ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  ErrorResponse "Invalid ID"
// @Failure      401  {object}  ErrorResponse "Unauthorized"
// @Failure      403  {object}  ErrorResponse "Forbidden"
// @Failure      500  {object}  ErrorResponse "Internal Server Error"
// @Router       /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
    // Парсим ID сообщения из URL
    vars := mux.Vars(r) // если используешь gorilla/mux
    messageID, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "invalid message id", http.StatusBadRequest)
        return
    }

    // Получаем юзера из контекста
    user, err := GetUserFromContext(r.Context())
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // Удаляем сообщение через usecase
    err = h.messageUseCase.DeleteMessage(r.Context(), messageID, user.ID)
    if err != nil {
        if strings.Contains(err.Error(), "unauthorized") {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}



func GetUserFromContext(ctx context.Context) (*models.User, error) {
    user, ok := ctx.Value(middleware.UserContextKey).(*models.User)
    if !ok || user == nil {
        return nil, errors.New("user not found in context")
    }
    return user, nil
}

type ErrorResponse struct {
    Message string `json:"message"`
}


// websockets

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
       return true
    },
}
func(h *MessageHandler) WsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
       fmt.Println("Error upgrading:", err)
       return
    }
    defer conn.Close()
    // Listen for incoming messages
    for {
       // Read message from the client
       _, message, err := conn.ReadMessage()
       if err != nil {
          fmt.Println("Error reading message:", err)
          break
       }
       fmt.Printf("Received: %s\\n", message)
       // Echo the message back to the client
       if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
          fmt.Println("Error writing message:", err)
          break
       }
    }
}
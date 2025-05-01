package handlers

import (
	"encoding/json"
	"net/http"
	"poshta/internal/usecase"
	"poshta/pkg/reqresp"
	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatService usecase.ChatService
}

func NewChatHandler(chatService usecase.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// @Summary Create a new chat
// @Description Create a new chat between two users
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param request body reqresp.CreateChatRequest true "Create chat request"
// @Success 201 {object} models.Chat "Chat created successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid request"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /chats [post]
func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req reqresp.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	chat, err := h.chatService.CreateChat(r.Context(), req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chat)
}


// @Summary Delete user chat
// @Description Delete user chat
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param chat_id path string true "Chat ID"
// @Success 200 {array} string "Chats deleted successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid chat ID"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /chats/{chat_id}/chats [delete]
func (h* ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	chatID := vars["chat_id"]

	if chatID == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	_, err := h.chatService.DeleteChat(r.Context(), chatID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chatID)

}

// @Summary Get user chats
// @Description Get all chats for a specific user
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Chat "Chats retrieved successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid user ID"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /chats/{user_id}/chats [get]
func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"] // no more Atoi

	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chats, err := h.chatService.GetUserChats(r.Context(), userID) // pass string now
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chats)
}


// @Summary Get chat messages
// @Description Get all messages for a specific chat
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param chat_id path string true "Chat ID"
// @Success 200 {array} models.Message "Messages retrieved successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid chat ID"
// @Failure 404 {object} reqresp.ErrorResponse "Chat not found"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /chats/{chat_id}/messages [get]
func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chat_id"] // no ParseInt anymore

	if chatID == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	// First check if chat exists
	chat, err := h.chatService.GetChatByID(r.Context(), chatID) // pass string now
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if chat == nil {
		respondWithError(w, http.StatusNotFound, "Chat not found")
		return
	}

	messages, err := h.chatService.GetChatMessages(r.Context(), chatID) // pass string
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, messages)
}


// Helper functions for responding with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, reqresp.ErrorResponse{
		Error: message,
	})
}


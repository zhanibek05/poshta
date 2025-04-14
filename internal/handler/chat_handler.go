package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"poshta/internal/service"
	"poshta/pkg/reqresp"
	"strconv"

	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
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

// @Summary Get user chats
// @Description Get all chats for a specific user
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Chat "Chats retrieved successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid user ID"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /chats/{user_id}/chats [get]
func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars, vars["user_id"])
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	chats, err := h.chatService.GetUserChats(r.Context(), (userID))
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
// @Param chat_id path int true "Chat ID"
// @Success 200 {array} models.Message "Messages retrieved successfully"
// @Failure 400 {object} reqresp.ErrorResponse "Invalid chat ID"
// @Failure 404 {object} reqresp.ErrorResponse "Chat not found"
// @Failure 500 {object} reqresp.ErrorResponse "Server error"
// @Router /chats/{chat_id}/messages [get]
func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID, err := strconv.ParseInt(vars["chat_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	// First check if chat exists
	chat, err := h.chatService.GetChatByID(r.Context(), chatID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if chat == nil {
		respondWithError(w, http.StatusNotFound, "Chat not found")
		return
	}

	messages, err := h.chatService.GetChatMessages(r.Context(), chatID)
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
package handlers


import (
	"encoding/json"
	"net/http"
	"poshta/internal/service"
	"poshta/pkg/reqresp"
)

type MessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
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

	message, err := h.messageService.SendMessage(r.Context(), req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, message)
}


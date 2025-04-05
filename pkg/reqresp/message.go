package reqresp

type SendMessageRequest struct {
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id"`
	Content  string `json:"content"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
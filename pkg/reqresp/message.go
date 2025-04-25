package reqresp

type SendMessageRequest struct {
	ChatID   int64  `json:"chat_id"`
	SenderID int64  `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content  string `json:"content"`
	EncryptedKey string `json:"encrypted_key"`
	EncryptedKeySender string `json:"encrypted_key_sender"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
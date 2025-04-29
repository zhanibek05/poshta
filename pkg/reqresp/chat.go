package reqresp

type CreateChatRequest struct {
	User1ID string `json:"user1_id"`
	User2ID string `json:"user2_id"`
}

type GetChatResponse struct{
	ChatID string  `json: "chat_id"`
	UserID string  `j son: "user_id"`
	Username string `json: "username"`
	PublicKey string  `json: "public_key`
}

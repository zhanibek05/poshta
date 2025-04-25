package reqresp

type CreateChatRequest struct {
	User1ID int64 `json:"user1_id"`
	User2ID int64 `json:"user2_id"`
}

type GetChatResponse struct{
	ChatID int64  `json: "chat_id"`
	UserID int64  `j son: "user_id"`
	Username string `json: "username"`
	PublicKey string  `json: "public_key`
}

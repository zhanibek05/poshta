package reqresp

type CreateChatRequest struct {
	Topic   string `json:"topic" binding:"required"`
	User1ID int64 `json:"user1_id"`
	User2ID int64 `json:"user2_id"`
}

package service

import (
	"context"
	"poshta/internal/domain/models"
	"poshta/internal/repository"
	"poshta/pkg/reqresp"
	"errors"
)

type MessageService interface {
	SendMessage(ctx context.Context, message reqresp.SendMessageRequest) (int64, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	userRepo 	repository.UserRepository
}

func NewMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository, userRepo repository.UserRepository) MessageService {
	return &messageService {
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		userRepo: 	 userRepo,	
	}
}

func (s *messageService) SendMessage(ctx context.Context, message reqresp.SendMessageRequest) (int64, error) {
	

	// Check if chat exists
	chat, err := s.chatRepo.GetByID(ctx, message.ChatID)
	if err != nil {
		return 0, err
	}
	if chat == nil {
		return 0, errors.New("chat not found")
	}
	// Create message

	// get username from user_id
	user, err := s.userRepo.GetByID(ctx, message.SenderID)



	messageModel := models.Message{
		ChatID:   message.ChatID,
		SenderID: message.SenderID,
		SenderName: user.Username,
		Content:  message.Content,
		EncryptedKey: message.EncryptedKey,
		EncryptedKeySender: message.EncryptedKeySender,
	}
	messageID, err := s.messageRepo.Create(ctx, &messageModel)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}
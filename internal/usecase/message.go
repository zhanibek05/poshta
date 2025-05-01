package usecase

import (
	"context"
	"errors"
	"fmt"
	"poshta/internal/domain/models"
	"poshta/internal/repository"
	"poshta/pkg/reqresp"
	"time"
)

type MessageUseCase interface {
	SendMessage(ctx context.Context, message reqresp.SendMessageRequest) (int64, error)
	DeleteMessage(ctx context.Context, messageID int64, requesterID string ) (error)
}

type messageUseCase struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	userRepo 	repository.UserRepository
}

func NewMessageUseCase(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository, userRepo repository.UserRepository) MessageUseCase {
	return &messageUseCase {
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		userRepo: 	 userRepo,	
	}
}

func (s *messageUseCase) SendMessage(ctx context.Context, message reqresp.SendMessageRequest) (int64, error) {
	

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
		CreatedAt: time.Now().UTC(),
	}
	messageID, err := s.messageRepo.Create(ctx, &messageModel)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}


func (u *messageUseCase) DeleteMessage(ctx context.Context, messageID int64, requesterID string) error {

	msg, err := u.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	if msg.SenderID != requesterID {
		return fmt.Errorf("unauthorized")
	}

	return u.messageRepo.Delete(ctx, messageID)
}
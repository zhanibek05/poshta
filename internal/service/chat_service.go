package service

import (
	"poshta/internal/domain/models"
	"context"
	"poshta/internal/repository"
	"poshta/pkg/reqresp"
)

type Chatservice interface {
	CreateChat(ctx context.Context, chat models.Chat) (models.Chat, error)
	GetUserChats(ctx context.Context, userID int) ([]models.Chat, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) chatService{
	return chatService{
		chatRepo: chatRepo,
	}
}

func (s chatService) CreateChat(ctx context.Context, req reqresp.CreateChatRequest) (models.Chat, error) {
	// check if chat exists
	existingChat, err := s.chatRepo.GetByUsersID(ctx, int64(req.User1ID), int64(req.User2ID))
	if err != nil {
		return models.Chat{}, err
	}
	if existingChat != nil {
		return *existingChat, nil
	}

	// create chat
	chat := models.Chat{
		User1ID: req.User1ID,
		User2ID: req.User2ID,
	}

	chatID, err := s.chatRepo.Create(ctx, &chat)
	if err != nil {
		return models.Chat{}, err
	}

	// get created chat 
	chatPtr, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return models.Chat{}, err
	}
	chat = *chatPtr
	return chat, nil
	
}

// get chats of users

func (s chatService) GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error) {
	// check if user exists
	existingUser, err := s.chatRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, nil
	}

	// get chats of user

	chats, err := s.chatRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return chats, nil

}
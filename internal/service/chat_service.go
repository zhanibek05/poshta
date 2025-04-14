package service

import (
	"poshta/internal/domain/models"
	"context"
	"poshta/internal/repository"
	"poshta/pkg/reqresp"
	"errors"
)

type ChatService interface {
	CreateChat(ctx context.Context, chat reqresp.CreateChatRequest) (models.Chat, error)
	GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error)
	GetChatByID(ctx context.Context, chatID int64) (*models.Chat, error)
	GetChatMessages(ctx context.Context, chatID int64) ([]models.Message, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
	userRepo repository.UserRepository
}

func NewChatService(chatRepo repository.ChatRepository, userRepo repository.UserRepository) ChatService {
    return &chatService{
        chatRepo: chatRepo,
		userRepo: userRepo,
    }
}


func (s *chatService) CreateChat(ctx context.Context, req reqresp.CreateChatRequest) (models.Chat, error) {
	existingChat, err := s.chatRepo.GetByUsersID(ctx, int64(req.User1ID), int64(req.User2ID))
	if err != nil {
		return models.Chat{}, err
	}
	if existingChat != nil {
		return *existingChat, nil
	}

	// Create chat
	chat := models.Chat{
		User1ID: req.User1ID,
		User2ID: req.User2ID,
	}
	chatID, err := s.chatRepo.Create(ctx, &chat)
	if err != nil {
		return models.Chat{}, err
	}

	// Get created chat
	chatPtr, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return models.Chat{}, err
	}
	if chatPtr == nil {
		return models.Chat{}, errors.New("failed to retrieve created chat")
	}
	
	return *chatPtr, nil
}
// get chats of users

func (s *chatService) GetUserChats(ctx context.Context, userID int64) ([]models.Chat, error) {
	// Check if user exists - using proper user repository
	// existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	// }
	// if existingUser != nil {
	// 	return nil, ErrUserExists
	// }

	// Get chats of user
	chats, err := s.chatRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *chatService) GetChatByID(ctx context.Context, chatID int64) (*models.Chat, error) {
	return s.chatRepo.GetByID(ctx, chatID)
}

func (s *chatService) GetChatMessages(ctx context.Context, chatID int64) ([]models.Message, error) {
	return s.chatRepo.GetMessages(ctx, chatID)
}
package usecase

import (
	"poshta/internal/domain/models"
	"context"
	"poshta/internal/repository"
	"poshta/pkg/reqresp"
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInternal           = errors.New("internal error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
)

type ChatService interface {
	CreateChat(ctx context.Context, chat reqresp.CreateChatRequest) (models.Chat, error)
	GetUserChats(ctx context.Context, userID string) ([]reqresp.GetChatResponse, error)
	GetChatByID(ctx context.Context, chatID string) (*models.Chat, error)
	GetChatMessages(ctx context.Context, chatID string, userID string) (reqresp.Chat, error)
	DeleteChat(ctx context.Context, chatID string) (string, error)
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

	// check users exist
	existingUser1, err := s.userRepo.GetByID(ctx, req.User1ID)
	if err != nil {	
		return models.Chat{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if existingUser1 == nil {
		return models.Chat{}, ErrUserNotFound	
	}

	existingUser2, err := s.userRepo.GetByID(ctx, req.User2ID)
	if err != nil {
		return models.Chat{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if existingUser2 == nil {
		return models.Chat{}, ErrUserNotFound
	}

	

	// check if chat with these users already exists
	existingChat, err := s.chatRepo.GetByUsersID(ctx, (req.User1ID), (req.User2ID))
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

func (s *chatService) GetUserChats(ctx context.Context, userID string) ([]reqresp.GetChatResponse, error) {
	
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	// Get chats of user
	
	chats, err := s.chatRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []reqresp.GetChatResponse
	for _, chat := range chats {
		// Determine the other user in the chat
		var otherUserID string
		if chat.User1ID == userID {
			otherUserID = chat.User2ID
		} else {
			otherUserID = chat.User1ID
		}

		// Fetch the other user's info
		otherUser, err := s.userRepo.GetByID(ctx, otherUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch other user: %w", err)
		}
		if otherUser == nil {
			continue // or return error if you want strict behavior
		}

		responses = append(responses, reqresp.GetChatResponse{
			ChatID:   chat.ID,
			UserID:   otherUser.ID,
			Username: otherUser.Username,
			PublicKey: otherUser.PublicKey,
		})
	}

	return responses, nil
}

func (s *chatService) GetChatByID(ctx context.Context, chatID string) (*models.Chat, error) {
	return s.chatRepo.GetByID(ctx, chatID)

}

func (s *chatService) GetChatMessages(ctx context.Context, chatID string, userID string) (reqresp.Chat, error) {
	messages, err := s.chatRepo.GetMessages(ctx, chatID)
	if err != nil {
		return reqresp.Chat{}, err
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return reqresp.Chat{}, err
	}

	var otherUserID string

	if chat.User1ID == userID {
		otherUserID = chat.User2ID
	} else {
		otherUserID = chat.User1ID
	}

	otherUser, err := s.userRepo.GetByID(ctx, otherUserID)
	if err != nil {
		return reqresp.Chat{}, err
	}

	return reqresp.Chat{
		ChatID:   chat.ID,
		Username: otherUser.Username,
		Messages: messages,	
	}, nil


	
}

func (s* chatService) DeleteChat(ctx context.Context, chatID string) (string, error) {
	// _, chats := s.GetUserChats(ctx, )
	return s.chatRepo.Delete(ctx, chatID)
}
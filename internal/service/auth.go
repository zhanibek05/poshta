package service

import (
	"context"
	"errors"
	"fmt"
	"poshta/internal/domain/models"
	"poshta/internal/repository"
	"time"
	"poshta/pkg/reqresp"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Auth-related errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInternal           = errors.New("internal error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
)


type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// AuthService defines the interface for authentication services
type AuthService interface {
	Register(ctx context.Context, req reqresp.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req reqresp.LoginRequest) (*reqresp.AuthResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserFromToken(token *jwt.Token) (*models.User, error)
	RefreshToken(refreshToken string) (*reqresp.AuthResponse, error)
	GetUserPublicKey(userID string) (string, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
	jwtCfg   JWTConfig
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.UserRepository, jwtCfg JWTConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// Register handles user registration
func (s *authService) Register(ctx context.Context, req reqresp.RegisterRequest) (*models.User, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}


	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		PublicKey: req.PublicKey,
	}

	// Create user
	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Get created user
	user, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

// Login handles user authentication
func (s *authService) Login(ctx context.Context, req reqresp.LoginRequest) (*reqresp.AuthResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, accessExpiry, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	refreshToken, _, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return &reqresp.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(accessExpiry).Seconds()),
		UserID:       user.ID,
	}, nil
}

// generateAccessToken creates a new JWT access token
func (s *authService) generateAccessToken(user *models.User) (string, time.Time, error) {
	expiryTime := time.Now().Add(s.jwtCfg.AccessTokenTTL)
	
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%s", user.ID),
		"username": user.Username,
		"email": user.Email,
		"exp":   expiryTime.Unix(),
		"iat":   time.Now().Unix(),
		"iss":   s.jwtCfg.Issuer,
		"type":  "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtCfg.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}

// generateRefreshToken creates a new JWT refresh token
func (s *authService) generateRefreshToken(user *models.User) (string, time.Time, error) {
	expiryTime := time.Now().Add(s.jwtCfg.RefreshTokenTTL)
	
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%s", user.ID),
		"exp":   expiryTime.Unix(),
		"iat":   time.Now().Unix(),
		"iss":   s.jwtCfg.Issuer,
		"type":  "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtCfg.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}

// ValidateToken validates a JWT token
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}

// GetUserFromToken extracts user information from a validated token
func (s *authService) GetUserFromToken(token *jwt.Token) (*models.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Get user ID from token
	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	

	user, err := s.userRepo.GetByID(context.Background(), userID)
	
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

func (s *authService) GetUserPublicKey(userID string) (string, error) {
	user, err := s.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	return user.PublicKey, nil
}

// RefreshToken handles token refresh
func (s *authService) RefreshToken(refreshToken string) (*reqresp.AuthResponse, error) {
	// Validate refresh token
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Check if token is refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	// Get user ID from token and fetch user
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	

	user, err := s.userRepo.GetByID(context.Background(), (userIDStr))
	
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Generate new tokens
	accessToken, accessExpiry, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, _, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &reqresp.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(accessExpiry).Seconds()),
	}, nil
}

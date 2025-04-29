package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"poshta/internal/domain/models"
	"poshta/internal/middleware"
	"poshta/internal/service"
	"poshta/pkg/logger"
	"poshta/pkg/reqresp"

	"github.com/gorilla/mux"
)

// AuthHandler handles auth-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}


// Register godoc
// @Summary Register new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param request body reqresp.RegisterRequest true "User registration data"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req reqresp.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", err, nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			http.Error(w, "User already exists", http.StatusConflict)
		default:
			logger.Error("Registration failed", err, nil)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles user authentication


// Login godoc
// @Summary Login user
// @Description Authenticate a user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body reqresp.LoginRequest true "User login credentials"
// @Success 200 {object} reqresp.AuthResponse "Authentication successful"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req reqresp.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", err, nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authResp, err := h.authService.Login(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		default:
			logger.Error("Login failed", err, nil)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResp)
}

// RefreshToken handles token refresh



// RefreshToken godoc
// @Summary Refresh tokens
// @Description Refresh access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body reqresp.RefreshTokenRequest true "Refresh token" 
// @Success 200 {object} reqresp.AuthResponse "Tokens refreshed successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid or expired token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type refreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body", err, nil)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authResp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.Error("Token refresh failed", err, nil)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResp)
}


// GetUser godoc
// @Summary Get user information
// @Description Example of protected route with JWT middleware
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{} "User profile information"
// @Failure 500 {object} map[string]string "Service unavailable"
// @Router /profile [get]
func (h *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Get user from context
    user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
    if !ok {
        http.Error(w, "User not found in context", http.StatusInternalServerError)
        return
    }
    
    // Create response with user data
    response := map[string]interface{}{
        "message": "Protected resource accessed",
        "user": map[string]interface{}{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,

        },
    }
    
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error creating response", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

// GetUser godoc
// @Summary Get user public key
// @Description Get user public key 
// @Tags user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]string "User public key retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Service unavailable"
// @Router /{user_id}/public_key [get]
func (h *AuthHandler) GetUserPublicKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	publicKey, err := h.authService.GetUserPublicKey(userID)
	if err != nil {
		logger.Error("Failed to get user's public key", err, nil)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"user_id":    fmt.Sprintf("%s", userID),
		"public_key": publicKey,
	})
}

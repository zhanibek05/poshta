package start

import (
	"fmt"
	"net/http"
	"poshta/internal/app/config"
	"poshta/internal/handler"
	"poshta/internal/middleware"
	"poshta/pkg/logger"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "poshta/docs"
)

func HTTP(cfg *config.Config, authHandler *handlers.AuthHandler, chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler , jwtMiddleware *middleware.JWTMiddleware) {
	// Initialize mux router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("status: It's all good man"))
		logger.Info("Health check requested", nil)
	}).Methods("GET")

	// Swagger docs
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Auth routes
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/auth/refresh", authHandler.RefreshToken).Methods("POST")
	// get user's public key
	router.HandleFunc("/api/{user_id}/public_key", authHandler.GetUserPublicKey).Methods("GET")

	// Chat routes
	router.Handle("/api/chats", jwtMiddleware.CreateAuthenticatedHandler(chatHandler.CreateChat)).Methods("POST")
	router.Handle("/api/chats/{user_id}/chats", jwtMiddleware.CreateAuthenticatedHandler( chatHandler.GetUserChats)).Methods("GET")
	router.Handle("/api/chats/{chat_id}/messages", jwtMiddleware.CreateAuthenticatedHandler(chatHandler.GetChatMessages)).Methods("GET")
	
	// Message routes
	router.Handle("/api/message", jwtMiddleware.CreateAuthenticatedHandler(messageHandler.SendMessage)).Methods("POST")
	router.Handle("/api/messages/{id}", jwtMiddleware.CreateAuthenticatedHandler(messageHandler.DeleteMessage)).Methods("DELETE")


	// Protected route example
	router.Handle("/api/profile", jwtMiddleware.CreateAuthenticatedHandler(authHandler.GetUserProfile)).Methods("GET")

	// websocket

	http.HandleFunc("/ws", messageHandler.WsHandler)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	logger.Info("Starting HTTP server", logrus.Fields{
		"address": addr,
		"swagger": fmt.Sprintf("http://%s/swagger/index.html", addr),
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(router)

	
	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		logger.Error("HTTP server failed", err, nil)
	}
}

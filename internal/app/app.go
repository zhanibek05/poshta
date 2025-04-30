package app

import (
	"poshta/internal/app/config"
	"poshta/internal/app/connections"
	"poshta/internal/app/start"
	"poshta/internal/app/ws"
	"poshta/internal/handler"
	"poshta/internal/middleware"
	"poshta/internal/repository"
	"poshta/internal/service"
	"poshta/internal/usecase"
	"poshta/pkg/logger"
)

func Run(configFiles ...string) {

	// Инициализация логгера
	logger.Init()
	logger.Info("Starting application", nil)

	// Загрузка конфигурации
	cfg, err := config.NewConfig(configFiles...)
	if err != nil {
		logger.Error("Failed to load config", err, nil)
		panic(err)
	}

	// Инициализация соединений
	conns, err := connections.NewConnections(cfg)
	if err != nil {
		logger.Error("Failed to initialize connections", err, nil)
		panic(err)
	}
	defer conns.Close()

	// init repos
	userRepo := repository.NewUserRepository(conns.DB)
	chatRepo := repository.NewChatRepository(conns.DB)
	messageRepo := repository.NewMessageRepository(conns.DB)

	// init services

	authService := service.NewAuthService(userRepo, service.JWTConfig{
		SecretKey:       cfg.JWT.SecretKey,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		Issuer:          cfg.JWT.Issuer,
	} )
	chatService := usecase.NewChatService(chatRepo, userRepo)
	messageService := usecase.NewMessageUseCase(messageRepo, chatRepo, userRepo)


	hub := ws.NewHub()
	go hub.Run()

	// init handlers
	authHandler := handlers.NewAuthHandler(authService)
	chatHandler := handlers.NewChatHandler(chatService)
	messageHandler := handlers.NewMessageHandler(messageService)
	

	wsHandler := handlers.NewWSHandler(hub, messageService, chatService)

	// init jwt middleware

	jwtMiddleware := middleware.NewJWTMiddleware(authService)

	// Запуск HTTP сервера
	logger.Info("Starting HTTP server", nil)
	start.HTTP(cfg, authHandler, chatHandler, messageHandler, wsHandler,  jwtMiddleware)
}

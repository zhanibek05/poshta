package app

import (
	"poshta/internal/app/config"
	"poshta/internal/app/connections"
	"poshta/internal/app/start"
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

	// Запуск HTTP сервера
	logger.Info("Starting HTTP server", nil)
	start.HTTP(cfg)
}

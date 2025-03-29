package start

import (
	"fmt"
	"poshta/internal/app/config"
	"poshta/pkg/logger"
	"net/http"
	"github.com/sirupsen/logrus"
)

func HTTP(cfg *config.Config) {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("status: Its all good man"))
		logger.Info("Health check requested susldjkf", nil)
	})
	

	addr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	logger.Info("Starting HTTP server", logrus.Fields{"address": addr})

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("HTTP server failed", err, nil)
	}
}

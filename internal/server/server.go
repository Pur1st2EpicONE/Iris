package server

import (
	"Iris/internal/config"
	"Iris/internal/logger"
	"Iris/internal/server/httpserver"
	"net/http"
)

type Server interface {
	Run() error
	Shutdown()
}

func NewServer(logger logger.Logger, config config.Server, handler http.Handler) Server {
	return httpserver.NewServer(logger, config, handler)
}

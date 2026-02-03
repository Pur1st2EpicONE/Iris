package app

import (
	"Iris/internal/cache"
	"Iris/internal/config"
	"Iris/internal/handler"
	"Iris/internal/logger"
	"Iris/internal/repository"
	"Iris/internal/server"
	"Iris/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/dbpg"
)

type App struct {
	logger  logger.Logger
	logFile *os.File
	server  server.Server
	ctx     context.Context
	cancel  context.CancelFunc
	cache   cache.Cache
	storage repository.Storage
}

func Boot() *App {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger, logFile := logger.NewLogger(config.Logger)

	db, err := connectDB(logger, config.Storage)
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err, "layer", "app")
	}

	cache, err := connectCache(logger, config.Cache)
	if err != nil {
		logger.LogFatal("app — failed to connect to cache", err, "layer", "app")
	}

	return wireApp(db, cache, logger, logFile, config)

}

func connectDB(logger logger.Logger, config config.Storage) (*dbpg.DB, error) {
	db, err := repository.ConnectDB(config)
	if err != nil {
		return nil, err
	}
	logger.LogInfo("app — connected to database", "layer", "app")
	return db, nil
}

func connectCache(logger logger.Logger, config config.Cache) (cache.Cache, error) {
	cache, err := cache.Connect(logger, config)
	if err != nil {
		return nil, err
	}
	logger.LogInfo("app — connected to cache", "layer", "app")
	return cache, nil
}

func wireApp(db *dbpg.DB, cache cache.Cache, logger logger.Logger, logFile *os.File, config config.Config) *App {

	ctx, cancel := newContext(logger)
	storge := repository.NewStorage(logger, config.Storage, db)
	service := service.NewService(logger, cache, storge)
	handler := handler.NewHandler(service)
	server := server.NewServer(logger, config.Server, handler)

	return &App{
		logger:  logger,
		logFile: logFile,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
		cache:   cache,
		storage: storge,
	}

}

func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		logger.LogInfo("app — received signal "+sig.String()+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()

	return ctx, cancel

}

func (a *App) Run() {

	go func() {
		if err := a.server.Run(); err != nil {
			a.logger.LogFatal("server run failed", err, "layer", "app")
		}
	}()

	<-a.ctx.Done()

	a.Stop()

}

func (a *App) Stop() {

	a.server.Shutdown()

	a.cache.Close()
	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}

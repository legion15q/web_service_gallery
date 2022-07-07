package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/internal/config"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		//logger.Error(err)
		return
	}

	// Dependencies
	mongoClient, err := Postgres.NewClient(cfg.Postgres.URI, cfg.Postgres.User, cfg.Postgres.Password)
	if err != nil {
		//logger.Error(err)

		return
	}

	db := mongoClient.Database(cfg.Postgres.Name)

	// Services, Repos & API Handlers
	tables := database.NewTables(db)
	services := service.NewServices(service.Deps{
		Tables:          tables,
		StorageProvider: storageProvider,
		Environment:     cfg.Environment,
		Domain:          cfg.HTTP.Host,
	})
	handlers := transport.NewHandler(services, tokenManager)

	services.Files.InitStorageUploaderWorkers(context.Background())

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			//logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	//logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		//logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		//logger.Error(err.Error())
	}
}

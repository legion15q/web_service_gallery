package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/internal/config"
	"web_app/internal/database"
	"web_app/internal/server"
	"web_app/internal/service"
	transport "web_app/internal/transport/http"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Error(err)
		return
	}

	// Dependencies
	//db, err := gorm.Open("postgres", connStr)

	connStr := "host = db port = 5432 user=admin password=root dbname=postgres sslmode=disable"
	db, err := gorm.Open("postgres", connStr)

	if err != nil {
		log.Errorf("cant connect to db:", err)
		return
	}
	//mongoClient, err := Postgres.NewClient(cfg.Postgres.URI, cfg.Postgres.User, cfg.Postgres.Password)
	//db := mongoClient.Database(cfg.Postgres.Name)

	// создаем таблицы (если не созданы)
	tables := database.NewTables(db)
	//создаем наши сервисы (некоторые на основе таблиц)
	services := service.NewServices(service.Deps{
		Tables: tables,
		//todo. Добавить TokenManager, PasswordHasher, Cacher
		File_storage_path: "images/",
		Environment:       cfg.Environment,
		Domain:            cfg.HTTP.Host,
	})

	handlers := transport.NewHandler(services)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server:", err)
			return
		}
	}()
	log.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: ", err)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err.Error())
		log.Error(err)
	}
}

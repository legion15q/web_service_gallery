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
	"github.com/sirupsen/logrus"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info(cfg)
	// Dependencies
	connStr := "host=" + cfg.Postgres.Host + " port=" + cfg.Postgres.Port + " user=" + cfg.Postgres.User +
		" password=" + cfg.Postgres.Password + " dbname=" + cfg.Postgres.Name + " sslmode=" + cfg.Postgres.Postgres_ssl_mode
	db, err := gorm.Open("postgres", connStr)

	if err != nil {
		logrus.Errorf("cant connect to db:", err)
		return
	}

	// создаем таблицы (если не созданы)
	tables := database.NewTables(db)
	//создаем наши сервисы (некоторые на основе таблиц)
	services := service.NewServices(service.Deps{
		Tables: tables,
		//todo. Добавить TokenManager, PasswordHasher, Cacher
		File_storage_path: cfg.FileStorage.Path_in_wm,
		Environment:       cfg.Environment,
		Domain:            cfg.HTTP.Host,
		TemplateFileName:  cfg.HTML.Templates.Picture_info,
	})

	handlers := transport.NewHandler(services)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("error occurred while running http server:", err)
			return
		}
	}()
	logrus.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logrus.Errorf("failed to stop server: ", err)
	}

	if err := db.Close(); err != nil {
		fmt.Println(err.Error())
		logrus.Error(err)
	}
}

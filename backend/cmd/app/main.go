package main

import (
	"github.com/kust1q/Zapp/backend/internal/config"
	"github.com/kust1q/Zapp/backend/internal/delivery/http"
	"github.com/kust1q/Zapp/backend/internal/repository"
	"github.com/kust1q/Zapp/backend/internal/repository/minio"
	"github.com/kust1q/Zapp/backend/internal/repository/postgres"
	"github.com/kust1q/Zapp/backend/internal/servers"
	"github.com/kust1q/Zapp/backend/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	cfg := config.Get()

	db, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		logrus.Fatalf("failed to init db: %s", err.Error())
	}

	mc, err := minio.NewMinioClient(cfg.Minio)
	if err != nil {
		logrus.Fatalf("failed to init minio client: %s", err.Error())
	}
	err = minio.CreateBucket(mc, cfg.Minio)
	if err != nil {
		logrus.Fatalf("failed to create minio bucket: %s", err.Error())
	}

	handler := http.NewHandler(
		service.NewAuthService(repository.NewAuthStorage(db, mc)),
	)

	srv := new(servers.Server)
	if err := srv.Run(cfg.App.Port, handler.InitRouters()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

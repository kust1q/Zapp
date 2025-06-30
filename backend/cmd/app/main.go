package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kust1q/Zapp/backend/internal/delivery/http"
	"github.com/kust1q/Zapp/backend/internal/repository"
	"github.com/kust1q/Zapp/backend/internal/servers"
	"github.com/kust1q/Zapp/backend/internal/service"
	"github.com/kust1q/Zapp/backend/pkg/repository/postgres"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	//logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error .env loading: %s", err.Error())
	}

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to init db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := http.NewHandler(services)

	srv := new(servers.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRouters()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

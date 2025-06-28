package main

import (
	"log"

	"github.com/kust1q/Zapp/backend/internal/delivery/http"
	"github.com/kust1q/Zapp/backend/internal/repository"
	"github.com/kust1q/Zapp/backend/internal/servers"
	"github.com/kust1q/Zapp/backend/internal/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handler := http.NewHandler(services)

	srv := new(servers.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRouters()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

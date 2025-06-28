package main

import (
	"log"

	"github.com/kust1q/Zapp/backend/internal/delivery/http"
	"github.com/kust1q/Zapp/backend/internal/repository"
	"github.com/kust1q/Zapp/backend/internal/servers"
	"github.com/kust1q/Zapp/backend/internal/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handler := http.NewHandler(services)

	srv := new(servers.Server)
	if err := srv.Run("8080", handler.InitRouters()); err != nil {
		log.Fatalf("error occured while running http server %s", err.Error())
	}
}

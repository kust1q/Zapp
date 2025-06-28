package main

import (
	"log"

	"github.com/kust1q/Zapp/backend/internal/delivery/http"
	"github.com/kust1q/Zapp/backend/internal/servers"
)

func main() {
	srv := new(servers.Server)
	handler := new(http.Handler)
	if err := srv.Run("8080", handler.InitRouters()); err != nil {
		log.Fatalf("error occured while running http server %s", err.Error())
	}
}

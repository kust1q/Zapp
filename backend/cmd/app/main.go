package main

import (
	"github.com/kust1q/Zapp/backend/internal/servers"
	"log"
)

func main() {
	srv := new(servers.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatalf("error occured while running http server %s", err.Error())
	}
}

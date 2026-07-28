package main

import (
	"gochat/step1_simple_echo/api"
	"gochat/step1_simple_echo/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	srv := api.New(api.Options{Addr: ":" + cfg.Port})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

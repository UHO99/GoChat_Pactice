package main

import (
	"gochat/config"
	"gochat/step5_simple_client_lifecycle/api"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	srv := api.New(api.Options{Addr: ":" + cfg.StepFivePort})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

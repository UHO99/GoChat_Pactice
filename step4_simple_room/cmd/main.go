package main

import (
	"gochat/config"
	"gochat/step4_simple_room/api"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	srv := api.New(api.Options{Addr: ":" + cfg.StepFourPort})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

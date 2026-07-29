package main

import (
	"gochat/config"
	"gochat/step6_room_persistDB/api"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	srv := api.New(api.Options{Addr: ":" + cfg.StepSixPort})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

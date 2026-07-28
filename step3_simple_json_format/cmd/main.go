package main

import (
	"gochat/config"
	"gochat/step3_simple_json_format/api"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	srv := api.New(api.Options{Addr: ":" + cfg.StepThreePort})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

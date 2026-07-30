package main

import (
	"context"
	"gochat/config"
	"gochat/db/store"
	"gochat/step6_room_persistDB/api"
	"gochat/step6_room_persistDB/servers/hub"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connection pgxpool : ", err)
	}
	defer pool.Close()

	h := hub.NewHub()
	st := store.New(pool)

	srv := api.New(h, st, api.Options{Addr: ":" + cfg.StepSixPort})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

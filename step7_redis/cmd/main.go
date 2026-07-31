package main

import (
	"context"
	"gochat/config"
	"gochat/db/store"
	"gochat/step7_redis/api"
	"gochat/step7_redis/servers/broker"
	"gochat/step7_redis/servers/hub"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("cannot connection Redis : ", err)
	}
	defer rdb.Close()

	b := broker.NewRedisBroker(rdb)
	h := hub.NewHub(b)
	st := store.New(pool)

	srv := api.New(h, st, api.Options{Addr: ":" + cfg.StepSevenPort})

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("cannot start server : ", err)
	}
}

package main

import (
	"DockerPostgreExample/cmd/serve/handler"
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/dbSettings"
	"DockerPostgreExample/internal/logger"
	"DockerPostgreExample/internal/redis"
	"DockerPostgreExample/internal/redisSettings"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	logger.Init()

	//postgres database start
	db, err := dbSettings.Initialize()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	defer db.Close()
	//postgres database end

	//redis database start
	redisConn := redisSettings.InitializeRedis()
	logger.Log.Info().Msgf("redis %v", redisConn.Ping(context.Background()))
	rdb := redis.NewRDB(redisConn)
	//redis database end

	subCh := rdb.Sub(context.Background(), "congrats")

	go func() {
		for msg := range subCh {
			logger.Log.Info().Msg(msg.Payload)
		}
	}()

	m := data.NewManager(db, rdb)
	h := handler.NewHandler(m)

	router := chi.NewRouter()
	h.Register(router)

	//err = http.ListenAndServe(":80", router)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("")
	}
}

// docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres
// docker run --name redis_example -p 6379:6379 --rm redis redis-server --requirepass "12345"

//docker exec -it redis_example /bin/sh         --зайти в редис контейнер
//redis-cli      								-- зайти в тулзу редиса
//auth 12345     								-- ввести пароль
//
//exit 											-- выйти из контейнера

package main

import (
	"DockerPostgreExample/cmd/serve/handler"
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/dbSettings"
	"DockerPostgreExample/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"net/http"
)

func main() {
	logger.Init()

	err := loggerStackShowcase()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}

	db, err := dbSettings.Initialize()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	defer db.Close()

	m := data.NewManager(db)
	h := handler.NewHandler(m)

	router := chi.NewRouter()
	h.Register(router)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("")
	}
}

//docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres

func loggerStackShowcase() error {
	err := errors.New("logger stack showcase")
	return err
}

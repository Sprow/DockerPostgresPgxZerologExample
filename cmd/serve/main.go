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
	log := logger.NewLogger()

	err := loggerStackShowcase()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	db, err := dbSettings.Initialize(log)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}
	defer db.Close()

	m := data.NewManager(db, log)
	h := handler.NewHandler(m, log)

	router := chi.NewRouter()
	h.Register(router)

	err = http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("")
	}
}

//docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres

func loggerStackShowcase() error {
	err := errors.New("logger stack showcase")
	return err
}

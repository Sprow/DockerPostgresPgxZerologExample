package main

import (
	"DockerPostgreExample/cmd/serve/handler"
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/dbSettings"
	"DockerPostgreExample/internal/logger"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	logger.Init()

	//err := loggerStackShowcase()
	//if err != nil {
	//	logger.Log.Error().Stack().Err(err).Msg("")
	//}

	db, err := dbSettings.Initialize()
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	defer db.Close()

	m := data.NewManager(db)
	h := handler.NewHandler(m)

	r := router.New()
	h.Register(r)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

//func loggerStackShowcase() error {
//	err := errors.New("logger stack showcase")
//	return err
//}

//docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres

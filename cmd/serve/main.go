package main

import (
	"DockerPostgreExample/cmd/serve/webserver"
	"DockerPostgreExample/internal/data"
	"DockerPostgreExample/internal/dbSettings"
	"DockerPostgreExample/internal/logger"
	"runtime"
)

func main() {
	runtime.MemProfileRate = 0
	runtime.GOMAXPROCS(16)
	log := logger.NewLogger()

	db, err := dbSettings.Initialize(log)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}
	defer db.Close()

	m := data.NewManager(db, log)
	cfg := webserver.Default()
	s := webserver.NewServer(cfg, log, m)

	err = s.Run()
	if err != nil {
		log.Panic().Msg("server start fail")
	}

	//err = http.ListenAndServe(":8080", router)
	//if err != nil {
	//	log.Fatal().Stack().Err(err).Msg("")
	//}
}

//docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres

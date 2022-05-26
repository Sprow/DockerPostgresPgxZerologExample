Example fasthttp beta 1.0

Docker + postgres + pgx.pool + zerolog + json.iterator 
Change net/http to fasthttp

without docker-compose

postgres

```docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres```

run main func in cmd/serve/main.go

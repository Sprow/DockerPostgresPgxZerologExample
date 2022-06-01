Example

Docker + postgres + pgx.pool + zerolog + json.iterator 

--------------------------------

run postgres db in docker
```docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres```

run server without docker 
run cmd/serve/main.go

----------------------------------------
second way to run

run both in docker containers

change POSTGRES_HOST to "database" (POSTGRES_HOST=database) in .env file and type 
```docker-compose build```

```docker-compose up```
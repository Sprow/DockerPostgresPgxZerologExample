Example

add redis 

postgres

```docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=Jonny -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=data_db --rm postgres```

redis

```docker run --name redis_example -p 6379:6379 --rm redis redis-server --requirepass "12345"```

run main() in cmd/main.go




OLD
```docker-compose build```

```docker-compose up```

To run nginx with upstream (server x2)

```docker-compose up --scale server=2```

go to => http://localhost:8080/
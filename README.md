Example

Docker + postgres + pgx.pool + zerolog + json.iterator +nginx

```docker-compose build```

```docker-compose up```

To run nginx with upstream (server x2)

```docker-compose up --scale server=2```

go to => http://localhost:8080/
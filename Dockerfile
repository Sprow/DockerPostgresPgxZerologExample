FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080 8080

RUN go build -o ./serve DockerPostgreExample/cmd/serve

CMD ["./serve"]


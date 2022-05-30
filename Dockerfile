FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 3000 3000

RUN go build -o ./serve DockerPostgreExample/cmd/serve

CMD ["./serve"]


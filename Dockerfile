FROM golang:1.18-alpine AS builder

RUN apk add git

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

#EXPOSE 80 80
EXPOSE 8080 8080

RUN go build -o ./serve DockerPostgreExample/cmd/serve

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app /app/

CMD ["./serve"]

#FROM golang:1.18-alpine AS builder
#
#RUN apk add git
#
#WORKDIR /go/src/app
#
#COPY go.mod go.sum ./
#
#RUN go mod download
#
#COPY . .
#
#EXPOSE 8080 8080
#
#RUN go build -o ./serve DockerPostgreExample/cmd/serve
#
#FROM alpine
#
#WORKDIR /app
#
#COPY --from=builder /go/src/app /app/
#
#CMD ["./serve"]
#

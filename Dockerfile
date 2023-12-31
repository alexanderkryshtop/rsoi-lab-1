FROM golang:1.20

WORKDIR /app
COPY src .
RUN go mod download

RUN go build -o app ./cmd

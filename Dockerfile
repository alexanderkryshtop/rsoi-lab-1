FROM golang:1.20

WORKDIR /app
COPY src .
RUN go mod download

RUN go build -o app ./cmd

ARG db_user=postgres
ENV POSTGRES_USER=$db_user

ARG db_pass=postgres
ENV POSTGRES_PASSWORD=$db_pass

ARG db_host=host.docker.internal
ENV POSTGRES_HOST=$db_host

ARG db_name=persons
ENV POSTGRES_DB=$db_name

CMD ["./app", "-c", "./cmd/config_local.yaml"]

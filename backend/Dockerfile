FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN go build -o bin/cloudmesh cmd/main.go

FROM alpine:3.21

RUN apk update

RUN apk add --no-cache make

RUN wget https://github.com/pressly/goose/releases/download/v3.24.1/goose_linux_x86_64 -O goose

RUN chmod +x goose

RUN mv goose /usr/local/bin/goose

COPY --from=build /src/bin/cloudmesh /bin/cloudmesh

COPY ./makefile .

COPY ./oauth.config.sample.yml .

COPY ./entrypoint.sh .

COPY ./setup_oauth_config.sh .

RUN mkdir -p sqlc/migrations

COPY ./sqlc/migrations/ sqlc/migrations

RUN chmod +x ./entrypoint.sh

RUN chmod +x ./setup_oauth_config.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]

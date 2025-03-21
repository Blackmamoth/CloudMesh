FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN go build -o bin/cloudmesh cmd/main.go

FROM alpine:3.21

RUN apk update

RUN apk add --no-cache make

COPY --from=build /src/bin/cloudmesh /bin/cloudmesh

COPY ./oauth.config.sample.yml .

COPY ./entrypoint.sh .

COPY ./setup_oauth_config.sh .

RUN chmod +x ./entrypoint.sh

RUN chmod +x ./setup_oauth_config.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]

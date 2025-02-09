FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN go build -o bin/cloudmesh cmd/main.go

FROM scratch

COPY --from=build /src/bin/cloudmesh /bin/cloudmesh

EXPOSE 8080

CMD ["/bin/cloudmesh"]

# dokkuapi
RESTful API for Dokku exposed on port `:3000`.

## Setup
Build from root of the repository: `go build -o dokkuapi main.go` or for linux deployment `GOOS=linux GOARCH=amd64 go build`\
Run without building `go run main.go`

## Endpoints
> `/info` temporary endpoint


## Logging
Loggin is provided by package [logger](github.com/ondro2208/dokkuapi/logger). It exposes `GeneralLogger` and `ErrorLogger` to log api runtime into `dokkuapi.log` file. Webserver traffice is logged sepparate to `dokkuapi_webserver.log` file via [gorilla/handler](https://github.com/gorilla/handlers).

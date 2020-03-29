# dokkuapi
RESTful API for Dokku exposed locally on port `:3000` or [server url](http://api.taron.sk).

## Setup
Build from root of the repository: `go build -o dokkuapi main.go` or for linux deployment `GOOS=linux GOARCH=amd64 go build`\
Run without building `go run main.go`

## Environmental variables
Don't forget to set all necessary ENVs. You can run it locally with `./start_with_env.sh`. On server it's required to set them in `/etc/systemd/system/dokkuapi.service` file.
* DB_URI
* DB_USERNAME
* DB_PWD
* JWT_TOKEN_SECRET


## Endpoints
> `/info` temporary endpoint\
> `/register` to register with github


## Logging
Logging is provided by package [logger](github.com/ondro2208/dokkuapi/logger). It exposes `GeneralLogger` and `ErrorLogger` to log api runtime into `dokkuapi.log` file. Webserver traffic is logged separate to `dokkuapi_webserver.log` file via [gorilla/handler](https://github.com/gorilla/handlers).

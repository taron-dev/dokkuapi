# dokkuapi
RESTful API for Dokku exposed locally on port `:3000` or on server http://api.taron.sk.

## Setup
Build from root of the repository: `go build -o dokkuapi main.go` or for linux deployment `GOOS=linux GOARCH=amd64 go build`\
Run without building `go run main.go`

## Environmental variables
Don't forget to set all necessary environmental variables. You can run it locally with `./start_with_env.sh`. On server it's required to set them in `/etc/systemd/system/dokkuapi.service` file. `EnvironmentFile` links to dokkuapi_env file, which contains also variables required for dokku, not only for dokkuapi.
* DB_URI
* DB_USERNAME
* DB_PWD
* JWT_TOKEN_SECRET

## Endpoints
|Path|Method|Description|
|----|------|-----------|
|/info|GET|temporary endpoint|
|/register|POST|to register with github|
|/login|POST|to login with github|
|/logout|POST|to deny access to api|
|/users/{userId}|DELETE|delete user according userId|
|/apps/|GET|retrieve info about user's applications|
|/apps/|POST|create application|
|/apps/{appId}|DELETE|delete user's application according appId|
|/apps/{appId}/services|POST|create and link backing service to application|
|/apps/{appId}/deploy|POST|use .tar file to deploy application|

## Tar deployment
Requires project directory compressed into `.tar` file. POST request should contains body form-data with `app_source_code` key and file value.


## Logging
Logging is provided by package [logger](github.com/ondro2208/dokkuapi/logger). It exposes `GeneralLogger` and `ErrorLogger` to log api runtime into `dokkuapi.log` file. Webserver traffic is logged separate to `dokkuapi_webserver.log` file via [gorilla/handler](https://github.com/gorilla/handlers).

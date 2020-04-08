package server

import (
	"github.com/gorilla/handlers"
	log "github.com/ondro2208/dokkuapi/logger"
	"os"
)

func (s *Server) initLogFile() {
	file, err := os.OpenFile("dokkuapi_webserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.ErrorLogger.Fatal(err)
	}
	defer file.Close()
	s.router = handlers.LoggingHandler(file, s.router)
}

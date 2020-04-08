package server

import (
	"encoding/json"
	log "github.com/ondro2208/dokkuapi/logger"
	"net/http"
)

func (s *Server) getInfo() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("content-type", "application/json")
		jsonData, err := json.Marshal(map[string]string{"message": "temporary info endpoint"})
		if err != nil {
			log.ErrorLogger.Fatal(err)
		}
		response.Write(jsonData)
	}
}

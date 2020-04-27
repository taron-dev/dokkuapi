package server

import (
	"github.com/ondro2208/dokkuapi/helper"
	"net/http"
)

func (s *Server) getInfo() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		helper.RespondWithMessage(response, request, http.StatusOK, "temporary info endpoint")
	}
}

package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

func ServicesGet(w http.ResponseWriter, r *http.Request, store *str.Store) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}

	services := []model.Service{}
	ss := service.NewServicesService(store)
	for _, serviceId := range app.Services {
		service, status, message := ss.GetService(serviceId)
		if service == nil {
			helper.RespondWithMessage(w, r, status, message)
		}
		services = append(services, *service)
	}
	helper.RespondWithData(w, r, http.StatusOK, services)

}

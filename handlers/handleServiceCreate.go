package handlers

import (
	"net/http"
	"strings"

	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/postgres"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
)

// ServiceCreate creates backing service for app
func ServiceCreate(w http.ResponseWriter, r *http.Request, store *str.Store) {
	//get appId param from url
	// authorize user towards appId
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}
	user, err := contextimpl.GetUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}

	//parse request body
	var servicePost postService
	err = helper.Decode(w, r, &servicePost)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, "Unable to parse request body")
		return
	}

	//dokku create service
	serviceType := strings.ToLower(servicePost.Type)

	var status int = -1
	var message string = ""
	switch serviceType {
	case "postgres":
		{
			status, message, err = postgres.Create(servicePost.Name, servicePost.Version)
			if err != nil {
				log.ErrorLogger.Println(err.Error())
				helper.RespondWithMessage(w, r, status, message)
				return
			}
		}
	default:
		{
			helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, "Unable to process service type")
			return
		}
	}

	//store service
	servicesService := service.NewServicesService(store)
	newService, status, message := servicesService.CreateService(servicePost.Name, serviceType)
	if newService == nil {
		//TODO destroy service in dokku
		log.ErrorLogger.Println("Store service - FAIL")
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	err = postgres.LinkServiceToApp(newService.Name, app.Name)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// update app services list
	usersService := service.NewUsersService(store)
	app.Services = append(app.Services, newService.Id)
	app, status, message = usersService.SetUserApplicationServices(*app, user.Id)
	if app == nil {
		log.ErrorLogger.Println("Update user application failed")
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	helper.RespondWithData(w, r, status, newService)
}

type postService struct {
	Name    string `json:"serviceName,omitempty"`
	Type    string `json:"serviceType,omitempty"`
	Version string `json:"serviceVersion,omitempty"`
}

package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/postgres"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
	"strings"
)

// ServiceCreate creates backing service for app
func ServiceCreate(w http.ResponseWriter, r *http.Request, store *str.Store) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := contextimpl.GetUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
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
	switch serviceType {
	case "postgres":
		{
			if serviceAlreadyExists := postgres.ServiceExists(servicePost.Name); serviceAlreadyExists {
				helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, "Wrong app name")
				return
			}
			ok, out := postgres.CreateService(servicePost.Name)
			if !ok {
				log.ErrorLogger.Println(out)
				helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Can't create postgres service", out)
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
		postgres.DestroyService(servicePost.Name)
		log.ErrorLogger.Println("Can't store service")
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	if ok, out := postgres.LinkService(newService.Name, app.Name); !ok {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Can't link service", out)
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

package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/postgres"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// ServiceDelete delete service
func ServiceDelete(w http.ResponseWriter, r *http.Request, store *str.Store) {
	user, err := contextimpl.GetUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// parse serviceID from url
	serviceIDParam, present := mux.Vars(r)["serviceId"]
	if !present || len(serviceIDParam) < 1 {
		log.ErrorLogger.Println("Extracting appId from url was unsuccessful")
		helper.RespondWithMessage(w, r, http.StatusNotFound, "Unknown service id")
		return
	}

	//find service name from DB
	ss := service.NewServicesService(store)
	serv, status, message := ss.GetServiceById(serviceIDParam)
	if status != http.StatusOK {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, message)
		return
	}
	if serv == nil {
		log.ErrorLogger.Println("Service with", serviceIDParam, "doesn't exists")
		helper.RespondWithMessage(w, r, http.StatusNotFound, "Unknown service")
		return
	}

	// unlink
	if ok, out := postgres.UnlinkService(serv.Name, app.Name); !ok {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Can't unlink service", out)
		return
	}

	// remove from app services
	us := service.NewUsersService(store)
	// prepare services ids array
	app.Services = removeStringFromArray(app.Services, serviceIDParam)
	app, status, message = us.SetUserApplicationServices(*app, user.Id)
	if app == nil {
		log.ErrorLogger.Println("Update user application failed")
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	//destroy
	if ok, out := postgres.DestroyService(serv.Name); !ok {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Can't destroy service", out)
		return
	}

	//remove from services collection
	err = ss.DeleteExistingService(serviceIDParam)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, "Service not deleted")
		return
	}
	helper.RespondWithMessage(w, r, http.StatusAccepted, "Service deleted")
}

func removeStringFromArray(s []primitive.ObjectID, r string) []primitive.ObjectID {
	for i, v := range s {
		if v.Hex() == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

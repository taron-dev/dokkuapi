package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/common"
	"github.com/ondro2208/dokkuapi/plugins/ps"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// AppsGet lists user's applications
func AppsGet(w http.ResponseWriter, r *http.Request, store *str.Store) {
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	log.GeneralLogger.Println("User id from jwt ", sub)

	userApps := []getApp{}

	//1. Get apps list from db
	usersService := service.NewUsersService(store)
	apps, status, message := usersService.GetUserApplications(sub)
	if status != http.StatusOK {
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	for _, app := range apps {
		userApp := new(getApp)
		userApp.ID = app.Id.Hex()

		userApp.Name = app.Name
		//2. read VHOST for each
		userApp.URLs = common.GetAppUrls(app.Name)
		// 3. GetRunningImageTag for each app
		userApp.Status = ps.GetAppStatus(app.Name)

		// 4. run scheduler ap status for each app
		if instances := common.GetAppInstances(app.Name); instances >= 0 {
			userApp.Instances = instances
		}

		// 5. read DOCKER_OPTIONS_DEPLOY file for each app
		if restartPolicy, err := ps.GetRestartPolicy(app.Name); err == nil {
			userApp.RestartPolicy = restartPolicy
		}

		userApps = append(userApps, *userApp)
	}
	helper.RespondWithData(w, r, http.StatusOK, userApps)
}

type getApp struct {
	ID            string   `json:"appId"`
	Name          string   `json:"appName"`
	URLs          []string `json:"urls"`
	Status        string   `json:"status"`
	Instances     int      `json:"instances"`
	RestartPolicy string   `json:"restartPolicy"`
}

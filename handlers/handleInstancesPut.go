package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/common"
	"github.com/ondro2208/dokkuapi/plugins/ps"
	"github.com/ondro2208/dokkuapi/service"
	"net/http"
)

// InstancesPut set new web and wokers instances count
func InstancesPut(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		log.ErrorLogger.Println("Can't get app from request's context:", err.Error())
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse put body
	putInstancesBody := new(putInstances)
	err = helper.Decode(w, r, putInstancesBody)
	if err != nil {
		log.ErrorLogger.Println("Can't decode request's body:", err.Error())
		helper.RespondWithMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Validate negative number
	if (putInstancesBody.WebCount < 0) || (putInstancesBody.WorkerCount < 0) || (putInstancesBody.WebCount < 1 && putInstancesBody.WorkerCount < 1) {
		log.ErrorLogger.Println("Wrong request's body:", putInstancesBody)
		helper.RespondWithMessage(w, r, http.StatusBadRequest, "Wrong request's body")
		return
	}

	// dokku ps:scale web=x worker=y
	err = ps.Scale(app.Name, putInstancesBody.WebCount, putInstancesBody.WorkerCount)
	if err != nil {
		log.ErrorLogger.Println("Can't scale app properly:", err.Error())
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Retreive instances information
	webContainerIDs, err := common.GetWebContainerIDs(app.Name)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusCreated, "App scaled successfully")
		return
	}
	is := service.NewInstancesService()
	instances, status, message := is.GetInstancesInfo(webContainerIDs)
	if instances == nil {
		helper.RespondWithMessage(w, r, status, message)
	}

	// Respond with actual information about instances
	helper.RespondWithData(w, r, http.StatusCreated, instances)
}

type putInstances struct {
	WebCount    int `json="webCount"`
	WorkerCount int `json="workerCount"`
}

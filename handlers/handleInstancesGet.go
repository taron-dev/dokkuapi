package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/common"
	"github.com/ondro2208/dokkuapi/service"

	"net/http"
)

// InstancesGet lists apps's instances
func InstancesGet(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	webContainerIDs, err := common.GetWebContainerIDs(app.Name)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	is := service.NewInstancesService()
	instances, status, message := is.GetInstancesInfo(webContainerIDs)
	if status != http.StatusOK {
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	helper.RespondWithData(w, r, status, instances)
}

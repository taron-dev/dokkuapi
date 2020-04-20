package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/common"
	"net/http"
)

// InstancesGet lists apps's instances
func InstancesGet(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}

	instances := []getInstance{}

	webContainerIDs, err := common.GetWebContainerIDs(app.Name)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}

	for _, containerID := range webContainerIDs {
		instance := new(getInstance)
		name, err := common.GetContainerName(containerID)
		if err != nil {
			helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		}
		instance.Name = name
		iType, err := common.GetContainerTypeFromName(name)
		if err != nil {
			helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		}
		instance.Type = iType
		status, err := common.GetContainerStatus(containerID)
		if err != nil {
			helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		}
		instance.Status = status

		instances = append(instances, *instance)
	}
	helper.RespondWithData(w, r, http.StatusOK, instances)
}

type getInstance struct {
	Name   string `json:"instanceName,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}

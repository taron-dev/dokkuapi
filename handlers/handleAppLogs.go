package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/logs"
	"net/http"
)

// AppLogs provide logs app
func AppLogs(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	logsOut, err := logs.GetAppLogs(app.Name, 20, "", false)
	if err != nil {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Cant' get logs", logsOut)
		return
	}
	logsResponse := new(getLogs)
	logsResponse.Logs = logsOut

	helper.RespondWithData(w, r, http.StatusOK, logsResponse)
}

type getLogs struct {
	Logs       string `json:"logs,omitempty"`
	FailedLogs string `json:"failedLogs,omitempty"`
}

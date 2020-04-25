package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/logs"
	"net/http"
)

// AppFailedLogs provide logs for last failed build
func AppFailedLogs(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	logsOut, err := logs.GetAppFailedLogs(app.Name)
	if err != nil {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Cant' get logs for failed build", logsOut)
		return
	}
	logsResponse := new(getLogs)
	logsResponse.FailedLogs = logsOut

	helper.RespondWithData(w, r, http.StatusOK, logsResponse)
}

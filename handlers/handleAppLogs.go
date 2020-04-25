package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/logs"
	"net/http"
	"net/url"
	"strconv"
)

// AppLogs provide logs app
func AppLogs(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	params := r.URL.Query()
	linesNum := parseIntQueryValue(params, "linesNum", -1)
	process := parseStringQueryValue(params, "process", "")
	quiet := parseBoolQueryValue(params, "quiet", false)

	logsOut, err := logs.GetAppLogs(app.Name, linesNum, process, quiet)
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

func parseStringQueryValue(params url.Values, key string, defaultValue string) string {
	keysMap, ok := params[key]
	if ok && (len(keysMap[0]) > 0) {
		return keysMap[0]
	}
	return defaultValue
}

func parseIntQueryValue(params url.Values, key string, defaultValue int) int {
	keysMap, ok := params[key]
	if ok && (len(keysMap[0]) > 0) {
		i, err := strconv.Atoi(keysMap[0])
		if err == nil {
			return i
		}
	}
	return defaultValue
}

func parseBoolQueryValue(params url.Values, key string, defaultValue bool) bool {
	keysMap, ok := params[key]
	if ok && (len(keysMap[0]) > 0) {
		val, err := strconv.ParseBool(keysMap[0])
		if err == nil {
			return val
		}
	}
	return defaultValue
}

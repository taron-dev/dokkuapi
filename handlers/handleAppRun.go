package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/run"
	"net/http"
)

//AppRun execute dokku run appName command
func AppRun(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// parse body
	body := new(runBody)
	err = helper.Decode(w, r, body)

	if err != nil || body.DokkuRun == "" {
		helper.RespondWithMessage(w, r, http.StatusBadRequest, "Can't process request body")
		return
	}

	// Call run
	output, err := run.DokkuRun(app.Name, body.DokkuRun)
	if err != nil {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Can't run dokku run command", output)
		return
	}
	helper.RespondWithMessageAndOutput(w, r, http.StatusCreated, "Command executed successfully", output)
}

type runBody struct {
	DokkuRun string `json:"dokkuRun,omitempty"`
}

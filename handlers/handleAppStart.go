package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/ps"
	"net/http"
)

func AppStart(w http.ResponseWriter, r *http.Request) {
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if ok, out := ps.StartApp(app.Name); !ok {
		helper.RespondWithMessageAndOutput(w, r, http.StatusInternalServerError, "Can't start app", out)
		return
	}
	helper.RespondWithMessage(w, r, http.StatusCreated, "App started successfully")
}

package handlers

import (
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
	"os/exec"
)

// AppEdit edit app fields
func AppEdit(w http.ResponseWriter, r *http.Request, store *str.Store) {
	// get app from conext
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := contextimpl.GetUser(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	usersService := service.NewUsersService(store)
	//parse body
	editApp := new(putApp)
	helper.Decode(w, r, editApp)
	if editApp.Name != "" {
		//dokku apps:rename appName newName
		// NOT WORKING - status, message, err := apps.RenameApp(app.Name, editApp.Name)
		out, err := exec.Command("dokku", "apps:rename", app.Name, editApp.Name).CombinedOutput()
		if err != nil {
			helper.RespondWithMessageAndOutput(w, r, http.StatusUnprocessableEntity, "Can't rename app", string(out))
			return
		}
		//update name in database
		status, message := usersService.SetUserApplicationName(user.Id, app.Id, editApp.Name)
		if status != http.StatusCreated {
			helper.RespondWithMessage(w, r, status, message)
			return
		}
	}

	editedApp := findApp(usersService, user.Id.Hex(), app.Id.Hex())
	if editedApp == nil {
		helper.RespondWithData(w, r, http.StatusCreated, editApp)
		return
	}
	helper.RespondWithData(w, r, http.StatusCreated, editedApp)
}

type putApp struct {
	Name string `json:"appName,omitempty"`
}

func findApp(us service.UsersService, userIDHex string, appIDHex string) *model.Application {
	userApps, status, _ := us.GetUserApplications(userIDHex)
	if status != http.StatusOK {
		return nil
	}
	for _, app := range userApps {
		if app.Id.Hex() == appIDHex {
			return &app
		}
	}
	return nil
}
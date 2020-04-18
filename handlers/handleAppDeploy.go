package handlers

import (
	"fmt"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/plugins/tar"
	"io/ioutil"
	"net/http"
	"os"
)

// AppDeploy provides deployment from tar file
func AppDeploy(w http.ResponseWriter, r *http.Request) {
	dokkuRoot := os.Getenv("DOKKU_ROOT")
	app, err := contextimpl.GetApp(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
	}

	r.ParseMultipartForm(25 << 20)
	requestFile, _, err := r.FormFile("app_source_code")
	if err != nil {
		log.ErrorLogger.Println("Error Retrieving the File")
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}
	defer requestFile.Close()
	fileDir := fmt.Sprintf("%v/%v", dokkuRoot, app.Name)

	fileBytes, err := ioutil.ReadAll(requestFile)
	if err != nil {
		log.ErrorLogger.Println("Error reading all")
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}
	appSourceCodeFilePath := fileDir + "/app_source_code.tar"
	err = ioutil.WriteFile(appSourceCodeFilePath, fileBytes, 0775)
	if err != nil {
		log.ErrorLogger.Println(err)
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	ok := tar.TarIn(app.Name, appSourceCodeFilePath)
	if !ok {
		log.ErrorLogger.Println("Error execute tar")
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, "Cannot deploy via tar")
		return
	}
	helper.RespondWithMessage(w, r, http.StatusCreated, "Successfully deployed")
}

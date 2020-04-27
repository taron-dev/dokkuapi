package handlers

import (
	"github.com/gorilla/mux"
	"github.com/ondro2208/dokkuapi/contextimpl"
	"github.com/ondro2208/dokkuapi/helper"
	"github.com/ondro2208/dokkuapi/plugins/ssh"
	"github.com/ondro2208/dokkuapi/service"
	str "github.com/ondro2208/dokkuapi/store"
	"net/http"
)

// UserEdit update user setup
func UserEdit(w http.ResponseWriter, r *http.Request, store *str.Store) {
	sub, err := contextimpl.GetSub(r.Context())
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	//if try to delete other user
	userIDParam := mux.Vars(r)["userId"]
	if sub != userIDParam {
		helper.RespondWithMessage(w, r, http.StatusUnauthorized, "Not Authorized")
		return
	}

	usersService := service.NewUsersService(store)
	user, status, message := usersService.GetExistingUserById(sub)
	if user == nil {
		helper.RespondWithMessage(w, r, status, message)
		return
	}

	// parse body
	body := new(editUserBody)
	err = helper.Decode(w, r, body)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusBadRequest, "Unable to parse request body")
		return
	}
	// validate ssh public key
	isValid, err := ssh.IsValidPublicSSHKey(user.Username, body.SSHPublicKey)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, "Unable to validate ssh public key")
		return
	}

	if !isValid {
		helper.RespondWithMessage(w, r, http.StatusUnprocessableEntity, "Ssh public key is not valid")
		return
	}

	// if ssh doesn't exists - add
	hasSSHKey, err := ssh.UserHasPublicSSHKey(user.Username)
	if err != nil {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, "Unable to retrieve info about ssh keys")
		return
	}
	if !hasSSHKey {
		if !ssh.AddSSHPublicKey(user.Username, body.SSHPublicKey) {
			helper.RespondWithMessage(w, r, http.StatusInternalServerError, "Can't add ssh public key")
			return
		}
	}

	// if already exists some key
	// remove old
	if !ssh.RemoveSSHPublicKey(user.Username) {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, "Can't update ssh public key")
		return
	}
	// add new
	if !ssh.AddSSHPublicKey(user.Username, body.SSHPublicKey) {
		helper.RespondWithMessage(w, r, http.StatusInternalServerError, "Old key was removed, but can't add new ssh public key")
		return
	}

	helper.RespondWithMessage(w, r, http.StatusCreated, "User updated successfully")
}

type editUserBody struct {
	SSHPublicKey string `json:"sshPublicKey"`
}

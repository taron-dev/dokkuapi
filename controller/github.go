package controller

import (
	"encoding/json"
	"errors"
	log "github.com/ondro2208/dokkuapi/logger"
	"net/http"
)

var httpClient = &http.Client{}

// GithubUser is model for github user
type GithubUser struct {
	Login string `json:"login,omitempty"`
	Id    int64  `json:"id,omitempty"`
}

const githubApiBaseUrl = "https://api.github.com"

// GetGithubUser fetch information about Github user
func GetGithubUser(accessToken string) (*GithubUser, error) {
	request, err := http.NewRequest("GET", githubApiBaseUrl+"/user", nil)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+accessToken)

	response, err := httpClient.Do(request)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("Not Authorized")
	}

	var githubUser *GithubUser
	err = json.NewDecoder(response.Body).Decode(&githubUser)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return githubUser, nil
}

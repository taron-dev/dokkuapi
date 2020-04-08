package authentication

import (
	"encoding/json"
	"errors"
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	"net/http"
)

var httpClient = &http.Client{}

const githubAPIBaseURL = "https://api.github.com"

// GetGithubUser fetch information about Github user
func GetGithubUser(accessToken string) (*model.GithubUser, error) {
	request, err := http.NewRequest("GET", githubAPIBaseURL+"/user", nil)
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

	var githubUser *model.GithubUser
	err = json.NewDecoder(response.Body).Decode(&githubUser)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return githubUser, nil
}

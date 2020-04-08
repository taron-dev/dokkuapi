package apps

import (
	dApps "github.com/dokku/dokku/plugins/apps"
	log "github.com/ondro2208/dokkuapi/logger"
)

// CreateApp : dokku apps:create appName
func CreateApp(appName string) (int, string, error) {
	err := dApps.CommandCreate([]string{appName})
	if err != nil {
		log.ErrorLogger.Println(err)
		return 422, "Can't create app", err
	}

	return 201, "", nil
}

// DestroyApp : dokku apps:destroy appName
func DestroyApp(appName string) (int, string, error) {
	err := dApps.CommandDestroy([]string{appName, "force"})
	if err != nil {
		log.ErrorLogger.Println(err)
		return 422, "Can't destory app", err
	}
	return 200, "", nil
}

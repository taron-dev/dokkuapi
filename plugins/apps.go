package plugins

import (
	"github.com/dokku/dokku/plugins/apps"
	log "github.com/ondro2208/dokkuapi/logger"
)

func CreateApp(appName string) (error, int, string) {
	err := apps.CommandCreate([]string{appName})
	if err != nil {
		log.ErrorLogger.Println(err)
		return err, 422, "Can't create app"
	}

	return nil, 201, ""
}

func DestroyApp(appName string) (error, int, string) {
	err := apps.CommandDestroy([]string{appName, "force"})
	if err != nil {
		log.ErrorLogger.Println(err)
		return err, 422, "Can't destory app"
	}
	return nil, 200, ""
}

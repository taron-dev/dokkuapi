package common

import (
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"os"
	"regexp"
	"strconv"
)

// GetAppUrls return app's urls
func GetAppUrls(appName string) []string {
	var dokkuRoot = os.Getenv("DOKKU_ROOT")
	path := fmt.Sprintf("%v/%v/VHOST", dokkuRoot, appName)
	urls, err := readUrls(path)
	if err != nil {
		log.ErrorLogger.Println("Can't read VHOST file from path:", path, "\n", err.Error())
		return []string{}
	}
	log.GeneralLogger.Println("VHOST read successfully")
	return urls
}

// GetAppStatus return if app is deployed or not
func GetAppStatus(appName string) string {
	if common.IsDeployed(appName) {
		return "DEPLOYED"
	}
	return "NOT DEPLOYED"
}

// GetAppInstances return app instances count
func GetAppInstances(appName string) int {
	scheduler := common.GetAppScheduler(appName)
	out, err := common.PlugnTriggerOutput("scheduler-app-status", []string{scheduler, appName}...)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return -1
	}
	re := regexp.MustCompile("\\d*")
	match := re.FindStringSubmatch(string(out))
	var value = 0
	if len(match) > 0 {
		stringVal := match[0]
		value, err = strconv.Atoi(stringVal)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return -1
		}
	}
	return value

}

package common

import (
	"errors"
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// GetWebContainerIDs provide web containers ID
func GetWebContainerIDs(appName string) ([]string, error) {
	return common.GetAppContainerIDs(appName, "web")
}

// GetWorkerContainerIDs provide worker containers ID
func GetWorkerContainerIDs(appName string) ([]string, error) {
	return common.GetAppContainerIDs(appName, "worker")
}

// GetContainerStatus get container status via docker inspect
func GetContainerStatus(containerID string) (string, error) {
	out, err := common.DockerInspect(containerID, "'{{.State.Status}}'")
	if err != nil {
		return "", err
	}
	return out, nil
}

// GetContainerName get container name via docker inspect
func GetContainerName(containerID string) (string, error) {
	out, err := common.DockerInspect(containerID, "'{{.Name}}'")
	if err != nil {
		return "", err
	}
	return strings.Split(out, "/")[1], nil
}

// GetContainerTypeFromName extract container type from name
func GetContainerTypeFromName(containerName string) (string, error) {
	r := regexp.MustCompile("\\w*\\.\\d")
	match := r.FindStringSubmatch(containerName)
	if len(match) > 0 {
		r = regexp.MustCompile("\\w*")
		match = r.FindStringSubmatch(match[0])
		if len(match) > 0 {
			return match[0], nil
		}
	}
	return "", errors.New("Can't parse container type from name")
}

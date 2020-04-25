package ps

import (
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"strings"
)

// Scale represents dokku ps:scale command
func Scale(appName string, webCount int, workerCount int) error {
	args := []string{"dokku", "ps:scale", appName}

	if webCount > 0 {
		webPart := fmt.Sprintf("web=%v", webCount)
		args = append(args, webPart)
	}

	if workerCount > 0 {
		workerPart := fmt.Sprintf("worker=%v", workerCount)
		args = append(args, workerPart)
	}

	log.GeneralLogger.Println(args)
	cmd := common.NewShellCmd(strings.Join(args, " "))
	cmd.ShowOutput = false
	out, err := cmd.Output()

	if err != nil {
		log.ErrorLogger.Println("Dokku ps:scale error:", err.Error())
		log.ErrorLogger.Println("Dokku ps:scale output:", string(out))
		return err
	}
	log.GeneralLogger.Println("Dokku ps:scale output:", string(out))
	return nil
}

// GetAppStatus return if app is deployed or not
func GetAppStatus(appName string) string {
	if !common.IsDeployed(appName) {
		return "NOT DEPLOYED"
	}

	webContainerIDs, err := common.GetAppContainerIDs(appName, "web")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return ""
	}
	workerContainerIDs, err := common.GetAppContainerIDs(appName, "worker")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return ""
	}

	isStopped := true

	for _, containerID := range webContainerIDs {
		status, err := common.DockerInspect(containerID, "'{{.State.Status}}'")
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return ""
		}
		if status != "exited" {
			isStopped = false
			break
		}
	}

	if !isStopped {
		return "DEPLOYED"
	}

	for _, containerID := range workerContainerIDs {
		status, err := common.DockerInspect(containerID, "'{{.State.Status}}'")
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return ""
		}
		if status != "exited" {
			isStopped = false
			break
		}
	}
	if !isStopped {
		return "DEPLOYED"
	}
	return "STOPPED"

}

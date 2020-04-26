package ps

import (
	"fmt"
	"github.com/dokku/dokku/plugins/apps"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"os/exec"
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

	err := apps.CommandLocked([]string{appName})
	if err == nil {
		return "BUILDING"
	}

	containerIDs, err := common.GetAppContainerIDs(appName, "")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return ""
	}
	for _, containerID := range containerIDs {
		status, err := common.DockerInspect(containerID, "'{{.State.Status}}'")
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return ""
		}
		if status != "exited" {
			return "DEPLOYED"
		}
	}

	return "STOPPED"
}

// StopApp dokku ps:stop appName
func StopApp(appName string) (bool, string) {
	out, err := exec.Command("dokku", "ps:stop", appName).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't stop app:", err.Error(), string(out))
		return false, string(out)
	}
	return true, ""
}

// StartApp dokku ps:start appName
func StartApp(appName string) (bool, string) {
	out, err := exec.Command("dokku", "ps:start", appName).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't start app:", err.Error(), string(out))
		return false, string(out)
	}
	return true, ""
}

// RestartApp dokku ps:restart appName
func RestartApp(appName string) (bool, string) {
	out, err := exec.Command("dokku", "ps:restart", appName).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't restart app:", err.Error(), string(out))
		return false, string(out)
	}
	return true, ""
}

// RebuildApp dokku ps:rebuild appName
func RebuildApp(appName string) (bool, string) {
	out, err := exec.Command("dokku", "ps:rebuild", appName).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't rebuild app:", err.Error(), string(out))
		return false, string(out)
	}
	return true, ""
}

package ps

import (
	"errors"
	"fmt"
	"github.com/dokku/dokku/plugins/apps"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
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

// GetRestartPolicy provides app's restart policy
func GetRestartPolicy(appName string) (string, error) {
	dokkuRoot := os.Getenv("DOKKU_ROOT")
	dockerOptionDeployPath := fmt.Sprintf("%v/%v/DOCKER_OPTIONS_DEPLOY", dokkuRoot, appName)
	contentBytes, err := ioutil.ReadFile(dockerOptionDeployPath)
	if err != nil {
		log.ErrorLogger.Println("Read DOCKER_OPTIONS_DEPLOY failed:", err.Error())
		return "", err
	}
	content := string(contentBytes)
	regex := regexp.MustCompile("--restart=.+")
	matches := regex.FindStringSubmatch(content)
	if len(matches) < 1 {
		return "", errors.New("Can't find restart policy")
	}
	regex = regexp.MustCompile("=.+")
	matches = regex.FindStringSubmatch(matches[0])
	if len(matches) < 1 {
		return "", errors.New("Can't find restart policy")
	}
	result := strings.Replace(matches[0], "=", "", -1)
	return result, nil
}

// SetRestartPolicy set app's restart policy
func SetRestartPolicy(appName string, policyVal string) error {
	dokkuRoot := os.Getenv("DOKKU_ROOT")
	dockerOptionDeployPath := fmt.Sprintf("%v/%v/DOCKER_OPTIONS_DEPLOY", dokkuRoot, appName)
	contentBytes, err := ioutil.ReadFile(dockerOptionDeployPath)
	if err != nil {
		log.ErrorLogger.Println("Read DOCKER_OPTIONS_DEPLOY failed:", err.Error())
		return err
	}
	content := string(contentBytes)
	regex := regexp.MustCompile("--restart=.+")
	matches := regex.FindStringSubmatch(content)
	if len(matches) < 1 {
		return errors.New("Can't read restart policy")
	}

	newContent := strings.Replace(content, matches[0], policyVal, -1)
	err = ioutil.WriteFile(dockerOptionDeployPath, []byte(newContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetValidPolicy provide valide policy value or empty string
func GetValidPolicy(policyName string, allowFailureCount int) string {
	if policyName == "on-failure" && allowFailureCount > 0 && allowFailureCount < 100 {
		return fmt.Sprintf("--restart=%v:%v", policyName, allowFailureCount)
	}

	policies := []string{"no", "unless-stopped", "always"}
	for _, policy := range policies {
		if policy == policyName {
			return fmt.Sprintf("--restart=%v", policy)
		}
	}
	return ""
}

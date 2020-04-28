package run

import (
	log "github.com/ondro2208/dokkuapi/logger"
	"os/exec"
	"regexp"
	"strings"
)

// DokkuRun represents dokku run appName command
func DokkuRun(appName string, cmdString string) (string, error) {
	args := []string{"run", appName}
	space := regexp.MustCompile(`\s+`)
	cmdCorrected := space.ReplaceAllString(cmdString, " ")
	cmdSlice := strings.Split(cmdCorrected, " ")
	args = append(args, cmdSlice...)
	out, err := exec.Command("dokku", args...).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't run 'dokku run':", err.Error(), string(out))
		return string(out), err
	}
	return string(out), nil
}

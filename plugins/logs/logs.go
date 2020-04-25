package logs

import (
	"fmt"
	log "github.com/ondro2208/dokkuapi/logger"
	"os/exec"
)

// GetAppLogs provides logs for application
func GetAppLogs(appName string, linesNum int, processName string, quiet bool) (string, error) {
	args := []string{"logs", appName}
	if linesNum > 0 {
		linesNumPart := fmt.Sprintf("-n %v", linesNum)
		args = append(args, linesNumPart)
	}

	if quiet {
		args = append(args, "-q")
	}

	if processName != "" {
		args = append(args, "-p")
		args = append(args, processName)
	}

	out, err := exec.Command("dokku", args...).CombinedOutput()
	output := string(out)
	log.GeneralLogger.Println("output:", output)
	log.GeneralLogger.Println("END")
	if err != nil {
		log.ErrorLogger.Println("Cant' execute logs command:", err.Error(), "\n", output)
		return output, err
	}
	return output, nil
}

// GetAppFailedLogs provides logs for last failed build
func GetAppFailedLogs(appName string) (string, error) {
	out, err := exec.Command("dokku", "logs:failed", appName).CombinedOutput()
	output := string(out)
	if err != nil {
		log.ErrorLogger.Println("Cant' execute failed logs command:", err.Error(), "\n", output)
		return output, err
	}
	return output, nil
}

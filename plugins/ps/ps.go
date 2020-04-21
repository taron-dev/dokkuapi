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

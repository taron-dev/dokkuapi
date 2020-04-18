package tar

import (
	"fmt"
	"github.com/dokku/dokku/plugins/common"
	log "github.com/ondro2208/dokkuapi/logger"
	"os"
	"os/exec"
	"strings"
)

// TarIn deploy file as dokku tar:in alternative
func TarIn(appName string, filePath string) bool {
	appPath := "/home/dokku/" + appName

	// Create source_code folder
	sourceCodeFolderPath := appPath + "/source_code"
	err := os.MkdirAll(sourceCodeFolderPath, 0775)
	if err != nil {
		log.ErrorLogger.Println("Can't create folder:", sourceCodeFolderPath)
		return false
	}

	//Untar file
	untarCmd := common.NewShellCmd(strings.Join([]string{"tar", "-xf", filePath, "--strip", "1", "-C", sourceCodeFolderPath}, " "))
	untarCmd.ShowOutput = false
	if ok := untarCmd.Execute(); !ok {
		log.ErrorLogger.Println("Can't untar loaded file:", filePath)
		return false
	}
	log.GeneralLogger.Println("Untared file:", filePath)

	// Deploy
	tarCmdString := fmt.Sprintf("tar c %v | dokku tar:in %v", sourceCodeFolderPath, appName)
	out, err := exec.Command("bash", "-c", tarCmdString).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		log.ErrorLogger.Println(string(out))
		return false
	}
	log.GeneralLogger.Println(appName, "deployed")
	return true
}

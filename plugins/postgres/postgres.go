package postgres

import (
	log "github.com/ondro2208/dokkuapi/logger"
	"os/exec"
)

// CreateService dokku postgres:create serviceName
func CreateService(serviceName string) (bool, string) {
	args := []string{"postgres:create", serviceName}
	out, err := exec.Command("dokku", args...).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't create service:", err.Error())
		return false, string(out)
	}
	log.GeneralLogger.Println("Create service output:", string(out))
	return true, string(out)
}

// DestroyService dokku postgres:destroy serviceName --force
func DestroyService(serviceName string) (bool, string) {
	out, err := exec.Command("dokku", "postgres:destroy", serviceName, "--force").CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't delete service:", err.Error())
		return false, string(out)
	}
	log.GeneralLogger.Println("Delete service output:", string(out))
	return true, string(out)
}

// LinkService dokku postgres:link serviceName appName
func LinkService(serviceName string, appName string) (bool, string) {
	out, err := exec.Command("dokku", "postgres:link", serviceName, appName).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't link service:", err.Error())
		return false, string(out)
	}
	log.GeneralLogger.Println("Link service output:", string(out))
	return true, string(out)
}

// UnlinkService dokku postgres:unlink serviceName appName
func UnlinkService(serviceName string, appName string) (bool, string) {
	out, err := exec.Command("dokku", "postgres:unlink", serviceName, appName).CombinedOutput()
	if err != nil {
		log.ErrorLogger.Println("Can't unlink service:", err.Error())
		return false, string(out)
	}
	log.GeneralLogger.Println("Unlink service output:", string(out))
	return true, string(out)
}

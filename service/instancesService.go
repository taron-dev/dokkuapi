package service

import (
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/model"
	"github.com/ondro2208/dokkuapi/plugins/common"
	"net/http"
)

// InstancesService interface methods
type InstancesService interface {
	GetInstancesInfo(containerIDs []string) ([]model.Instance, int, string)
}

// NewInstancesService Constructor
func NewInstancesService() InstancesService {
	return &InstancesServiceContext{}
}

// InstancesServiceContext struct
type InstancesServiceContext struct {
}

// GetInstancesInfo returns array with instances information
func (is *InstancesServiceContext) GetInstancesInfo(containerIDs []string) ([]model.Instance, int, string) {
	instances := []model.Instance{}
	for _, containerID := range containerIDs {
		instance := new(model.Instance)
		name, err := common.GetContainerName(containerID)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return nil, http.StatusInternalServerError, "Can't retrieve instance's container name"
		}
		instance.Name = name
		iType, err := common.GetContainerTypeFromName(name)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return nil, http.StatusInternalServerError, "Can't retrieve instance's container type"
		}
		instance.Type = iType
		status, err := common.GetContainerStatus(containerID)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return nil, http.StatusInternalServerError, "Can't retrieve instance's container status"
		}
		instance.Status = status

		instances = append(instances, *instance)
	}
	return instances, http.StatusOK, ""
}

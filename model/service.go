package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service is model for backing service in dokku
type Service struct {
	Id   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"serviceName,omitempty" bson:"serviceName,omitempty"`
	Type string             `json:"serviceType,omitempty" bson:"serviceType,omitempty"`
}

// type ServiceType string
// const (
// 	REDIS    ServiceType = "REDIS"
// 	POSTGRES ServiceType = "POSTGRES"
// )

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

// ServicePost represents expected body for POST request on /apps/{appId}/services endpoint
type ServicePost struct {
	Name    string `json:"serviceName,omitempty"`
	Type    string `json:"serviceType,omitempty"`
	Version string `json:"serviceVersion,omitempty"`
}

// type ServiceType string
// const (
// 	REDIS    ServiceType = "REDIS"
// 	POSTGRES ServiceType = "POSTGRES"
// )

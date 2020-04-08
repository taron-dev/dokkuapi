package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Application is model for backing service in dokku
type Application struct {
	Id       primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string               `json:"appName,omitempty" bson:"appName,omitempty"`
	Services []primitive.ObjectID `json:"services,omitempty" bson:"services,omitempty"`
}

// ApplicationPost represents expected body for POST request on /apps endpoint
type ApplicationPost struct {
	Name string `json:"appName,omitempty" bson:"appName,omitempty"`
}

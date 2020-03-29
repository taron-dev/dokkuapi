package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is model for backing service in dokku
type User struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username,omitempty" bson:"username,omitempty"`
	GithubId     int64              `json:"githubId,omitempty" bson:"githubId,omitempty"`
	Applications []Application      `json:"applications,omitempty" bson:"applications,omitempty"`
}

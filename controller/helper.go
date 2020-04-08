package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var client *mongo.Client

//GetCollection provide collection and context according input parameter
func GetCollection(client *mongo.Client, dbName string, collectionName string) (*mongo.Collection, context.Context) {
	client = client
	collection := client.Database(dbName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return collection, ctx
}

//TODO remove
func GetCollection1(collectionName string) (*mongo.Collection, context.Context) {
	collection := client.Database("").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return collection, ctx
}

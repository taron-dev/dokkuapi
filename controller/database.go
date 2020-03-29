package controller

import (
	"context"
	log "github.com/ondro2208/dokkuapi/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

// Client represents mongodb client
var client *mongo.Client = initializeDatabase()

const dbName = "dokkuapi"

// InitializeDatabase setup connection do dokkuapi database
func initializeDatabase() *mongo.Client {
	username := os.Getenv("DB_USERNAME")
	pwd := os.Getenv("DB_PWD")
	dbUri := os.Getenv("DB_URI")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().
		ApplyURI(dbUri).
		SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    dbName,
			Username:      username,
			Password:      pwd,
		})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.ErrorLogger.Fatal("Could not connect to MongoDB:", err)
	}
	return client
}

// GetCollection provide collection and context according input parameter
func GetCollection(collectionName string) (*mongo.Collection, context.Context) {
	collection := client.Database(dbName).Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return collection, ctx
}

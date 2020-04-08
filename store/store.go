package store

import (
	"context"
	log "github.com/ondro2208/dokkuapi/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const dbName = "dokkuapi"

type Store struct {
	Client *mongo.Client
	DbName string
}

func NewStore() (*Store, error) {
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

	str := new(Store)
	str.Client = client
	str.DbName = dbName

	return str, nil
}

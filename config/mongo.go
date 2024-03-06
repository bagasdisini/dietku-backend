package config

import (
	"context"
	"dietku-backend/cmd/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var db *mongo.Database

func ConnectMongo() *mongo.Database {
	if db == nil {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			log.Fatal("You must set your 'MONGODB_NMM_URI' environmental variable.")
		}

		dbName := os.Getenv("MONGODB_NAME")
		if dbName == "" {
			log.Fatal("You must set your 'MONGODB_NMM_NAME' environmental variable.")
		}

		// Set client options
		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			panic(err)
		}

		// Get database
		db = client.Database(dbName)
	}

	return db
}

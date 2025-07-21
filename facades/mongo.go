package facades

import (
	"golang_starter_kit_2025/config"

	"go.mongodb.org/mongo-driver/mongo"
)

func MongoClient() *mongo.Client {
	if config.MongoClient == nil {
		panic("MongoDB client is not initialized. Call config.ConnectMongo() first.")
	}
	return config.MongoClient
}

func MongoDatabase(name ...string) *mongo.Database {
	cfg := config.GetMongoDBConfig()
	dbName := cfg.Database
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}
	return MongoClient().Database(dbName)
}

func MongoCollection(collectionName string, dbName ...string) *mongo.Collection {
	if collectionName == "" {
		panic("Collection name must not be empty")
	}
	cfg := config.GetMongoDBConfig()
	databaseName := cfg.Database
	if len(dbName) > 0 && dbName[0] != "" {
		databaseName = dbName[0]
	}
	return MongoClient().Database(databaseName).Collection(collectionName)
}

package config

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Mongo() (*mongo.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	uri := "mongodb://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@" + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT") + "/default?authSource=" + os.Getenv("MONGO_AUTH_SOURCE")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	MongoClient = client
	return client, nil
}

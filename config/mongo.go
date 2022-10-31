package config

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	AuthSource string
}

func InitMongo() *MongoConfig {
	err := godotenv.Load(".env")
	if err != nil {
		return nil
	}
	return &MongoConfig{
		Host:       os.Getenv("MONGO_HOST"),
		Port:       os.Getenv("MONGO_PORT"),
		Username:   os.Getenv("MONGO_USERNAME"),
		Password:   os.Getenv("MONGO_PASSWORD"),
		AuthSource: os.Getenv("MONGO_AUTH_SOURCE"),
	}
}

func (mongoConfig *MongoConfig) Mongo() (*mongo.Client, error) {
	uri := "mongodb://" + mongoConfig.Username + ":" + mongoConfig.Password + "@" + mongoConfig.Host + ":" + mongoConfig.Port + "/default?authSource=" + mongoConfig.AuthSource
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}

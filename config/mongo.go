package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfig struct {
	Host              string
	Port              string
	Username          string
	Password          string
	AuthSource        string
	Db                string
	AccountCollection string
}

func InitMongo() *MongoConfig {
	err := godotenv.Load(".env")
	if err != nil {
		return nil
	}
	return &MongoConfig{
		Host:              os.Getenv("MONGO_HOST"),
		Port:              os.Getenv("MONGO_PORT"),
		Username:          os.Getenv("MONGO_USERNAME"),
		Password:          os.Getenv("MONGO_PASSWORD"),
		AuthSource:        os.Getenv("MONGO_AUTH_SOURCE"),
		Db:                os.Getenv("MONGO_DB"),
		AccountCollection: os.Getenv("MONGO_ACCOUNT_COLLECTION"),
	}
}

func (mongoConfig *MongoConfig) Mongo() (*mongo.Client, error) {
	// uri := "mongodb://" + mongoConfig.Username + ":" + mongoConfig.Password + "@" + mongoConfig.Host + ":" + mongoConfig.Port + "/default?authSource=" + mongoConfig.AuthSource
	uri := "mongodb://" + mongoConfig.Host + ":" + mongoConfig.Port
	credential := options.Credential{
		// AuthSource:    os.Getenv("MONGO_AUTH_SOURCE"),
		// AuthMechanism: "SCRAM-SHA-1",
		Username:      os.Getenv("MONGO_USERNAME"),
		Password:      os.Getenv("MONGO_PASSWORD"),
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())
	return client, nil
}

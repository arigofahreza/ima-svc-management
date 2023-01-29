package config

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client

func Redis() (*redis.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	intPort, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:       intPort,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	RedisClient = client
	return client, nil
}

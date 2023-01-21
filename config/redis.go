package config

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client

func Redis() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	intPort, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return err
	}
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   intPort,
	})
	defer client.Close()
	RedisClient = client
	return nil
}

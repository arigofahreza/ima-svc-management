package config

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Host     string
	Port     string
	Db       string
}

func InitRedis() *RedisConfig {
	err := godotenv.Load(".env")
	if err != nil {
		return nil
	}
	return &RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Db:       os.Getenv("REDIS_DB"),
	}
}

func (redisConfig RedisConfig) Redis() (*redis.Client, error) {
	intPort, err := strconv.Atoi(redisConfig.Db)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		DB:       intPort,
	})
	return client, nil
}

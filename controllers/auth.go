package controllers

import (
	"context"
	"errors"
	"ima-svc-management/config"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
)

var ACCESS_TOKEN = time.Duration(5) * time.Minute
var REFRESH_TOKEN = time.Duration(10) * time.Minute

type AuthController struct{}

func (authController AuthController) Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	redisClient, err := config.InitRedis().Redis()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("account")

	filter := bson.M{"name": username}

	account := model.AccountModel{}

	err = collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if compare, _ := helpers.PasswordCompare([]byte(password), []byte(account.Password)); !compare && err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong Password"})
		c.Abort()
		return
	}

	accessToken, err := helpers.GenerateToken(account.Id + ACCESS_TOKEN.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	refreshToken, err := helpers.GenerateToken(account.Id + ACCESS_TOKEN.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	getToken := redisClient.Get(context.Background(), accessToken)
	if errors.Is(err, redis.Nil) {
		setAccessToken := redisClient.Set(context.Background(), accessToken, account, ACCESS_TOKEN)
		if err := setAccessToken.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		setRefreshToken := redisClient.Set(context.Background(), refreshToken, account, REFRESH_TOKEN)
		if err := setRefreshToken.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	session.Se

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

func (authController AccountController) Logout(c *gin.Context) {

}

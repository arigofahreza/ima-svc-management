package controllers

import (
	"context"
	"ima-svc-management/config"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ACCESS_TOKEN = time.Duration(5) * time.Minute
var REFRESH_TOKEN = time.Duration(10) * time.Minute

type AuthController struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
}

func InitAuth() *AuthController {
	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		return nil
	}
	redisClient, err := config.InitRedis().Redis()
	if err != nil {
		return nil
	}
	return &AuthController{
		MongoClient: mongoClient,
		RedisClient: redisClient,
	}

}

// @Summary Login
// @Description Login user
// @Param email formData string true "email"
// @Param password formData string true "password"
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/auth/login [post]
func (authController AuthController) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	collection := authController.MongoClient.Database("test").Collection("account")

	filter := bson.M{"email": email}

	account := model.AccountModel{}

	err := collection.FindOne(context.TODO(), filter).Decode(&account)
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

	getToken, err := authController.RedisClient.Get(context.Background(), accessToken).Result()
	if getToken == "" {
		setAccessToken := authController.RedisClient.Set(context.Background(), accessToken, "", ACCESS_TOKEN)
		if err := setAccessToken.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		setRefreshToken := authController.RedisClient.Set(context.Background(), refreshToken, "", REFRESH_TOKEN)
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

	c.Header("Authorization", "Bearer " + getToken)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login Success"})
}

// @Summary Logout
// @Description Logout user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/auth/logout [post]
func (authController AuthController) Logout(c *gin.Context) {
	token := c.GetHeader("Token")

	authController.RedisClient.Del(context.Background(), token)

	c.Header("Token", "")
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Logout Success"})
}

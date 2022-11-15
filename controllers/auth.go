package controllers

import (
	"context"
	"ima-svc-management/config"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var ACCESS_TOKEN = time.Duration(5) * time.Minute
var REFRESH_TOKEN = time.Duration(10) * time.Minute

type AuthController struct{}

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

	filter := bson.M{"email": email}

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

	getToken, err := redisClient.Get(context.Background(), accessToken).Result()
	if getToken == "" {
		setAccessToken := redisClient.Set(context.Background(), accessToken, "", ACCESS_TOKEN)
		if err := setAccessToken.Err(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		setRefreshToken := redisClient.Set(context.Background(), refreshToken, "", REFRESH_TOKEN)
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

	c.Header("Token", getToken)

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

	redisClient, err := config.InitRedis().Redis()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	redisClient.Del(context.Background(), token)

	c.Header("Token", "")
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Logout Success"})
}

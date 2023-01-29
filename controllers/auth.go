package controllers

import (
	"context"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"

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
	Auth        *helpers.Auth
}

func InitAuth(redisClient *redis.Client, mongoClient *mongo.Client) *AuthController {
	return &AuthController{
		MongoClient: mongoClient,
		RedisClient: redisClient,
		Auth:        &helpers.Auth{},
	}
}

// @Summary Login
// @Description Login user
// @Param body body model.LoginModel true "body"
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/auth/login [post]
func (authController AuthController) Login(c *gin.Context) {
	ctx := context.Background()
	login := model.LoginModel{}
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := authController.MongoClient.Database("test").Collection("account")

	filter := bson.M{"email": login.Email}

	account := model.AccountModel{}

	err = collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email are not registered"})
			c.Abort()
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if compare, _ := helpers.PasswordCompare([]byte(login.Password), []byte(account.Password)); !compare && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong Password"})
		c.Abort()
		return
	}

	tokenDetails, err := authController.Auth.CreateToken(login.Email)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Failed creating token"})
		c.Abort()
		return
	}

	err = authController.Auth.CreateAuth(ctx, login.Email, tokenDetails)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	c.Header("Authorization", "Bearer "+tokenDetails.AccessToken)
	c.SetCookie("refresh_token", tokenDetails.RefreshToken, 86400, "/", "localhost", false, true)

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
	auth, err := authController.Auth.ExtractTokenMetadata(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	deleted, err := authController.Auth.DeleteAuth(context.Background(), auth.AccessUUID)
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Logout Success"})
}

// @Summary Refresh
// @Description refresh token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/auth/refresh [get]
func (authController AuthController) Refresh(c *gin.Context) {
	ctx := context.Background()
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := authController.Auth.VerifyRefreshToken(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token expired"})
		c.Abort()
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		email := claims["email"].(string)
		deleted, err := authController.Auth.DeleteAuth(ctx, refreshUuid)
		if err != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		newToken, err := authController.Auth.CreateToken(email)
		if err != nil {
			c.JSON(http.StatusForbidden, err.Error())
			return
		}

		err = authController.Auth.CreateAuth(ctx, email, newToken)
		if err != nil {
			c.JSON(http.StatusForbidden, err.Error())
			return
		}

		c.Header("Authorization", "Bearer "+newToken.AccessToken)
		c.SetCookie("refresh_token", newToken.RefreshToken, 86400, "/", "localhost", false, true)

	}

}

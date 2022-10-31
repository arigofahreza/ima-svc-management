package controllers

import (
	"context"
	"encoding/hex"
	"ima-svc-management/config"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountController struct{}

func (accountController AccountController) AddAccount(c *gin.Context) {

	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("account")

	account := model.AccountModel{}
	err = c.BindJSON(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	dataAccount := bson.M{
		"name":       account.Name,
		"email":      account.Email,
		"password":   helpers.GeneratePasswordHash([]byte(account.Password)),
		"created_at": time.Now(),
		"updated_at": nil,
	}

	hashId, err := bson.Marshal(dataAccount)
	if err != nil {
		log.Fatal(err)
	}

	dataAccount["id"] = hex.EncodeToString(hashId)

	_, err = collection.InsertOne(context.Background(), dataAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Create account successful"})

}

func (accountController AccountController) GetAccount(c *gin.Context) {

}

func (accountController AccountController) GetAccountById(c *gin.Context) {

}

func (accountController AccountController) UpdateAccount(c *gin.Context) {

}

func (accountController AccountController) DeleteAccount(c *gin.Context) {

}

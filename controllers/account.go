package controllers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"ima-svc-management/config"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"name":      account.Name,
		"email":     account.Email,
		"password":  helpers.GeneratePasswordHash([]byte(account.Password)),
		"createdAt": time.Now().Unix(),
		"updatedAt": nil,
	}

	hashId, err := bson.Marshal(dataAccount)
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.Sum(hashId)

	dataAccount["_id"] = hex.EncodeToString(hash[:])

	_, err = collection.InsertOne(context.Background(), dataAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Create account successful"})

}

func (accountController AccountController) GetAccount(c *gin.Context) {

	paginationModel := model.PaginateAccountModel{}

	err := c.BindJSON(&paginationModel)
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

	pageOptions := options.Find()
	pageOptions.SetLimit(int64(paginationModel.Size))
	if paginationModel.Order != "" && paginationModel.OrderBy != "" {
		if paginationModel.Order == "asc" {
			pageOptions.SetSort(bson.M{paginationModel.OrderBy: 1})
		} else {
			pageOptions.SetSort(bson.M{paginationModel.OrderBy: -1})
		}
	}
	pageOptions.SetSkip(int64(paginationModel.Page))

	cursor, err := collection.Find(context.TODO(), bson.D{{}}, pageOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	datas := make([]map[string]interface{}, 0)
	for cursor.Next(context.TODO()) {
		account := model.AccountModel{}
		if err := cursor.Decode(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		data := map[string]interface{}{
			"id":         account.Id,
			"name":       account.Name,
			"email":      account.Email,
			"created_at": account.CreatedAt,
			"updated_at": account.UpdatedAt,
		}
		datas = append(datas, data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})

}

func (accountController AccountController) GetAccountById(c *gin.Context) {

	account := model.AccountModel{}

	id := c.Query("id")

	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("account")

	filter := bson.M{"_id": id}

	err = collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	datas := make([]map[string]interface{}, 0)
	data := map[string]interface{}{
		"id":         account.Id,
		"name":       account.Name,
		"email":      account.Email,
		"created_at": account.CreatedAt,
		"updated_at": account.UpdatedAt,
	}
	datas = append(datas, data)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

func (accountController AccountController) UpdateAccount(c *gin.Context) {

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

	filter := bson.M{"_id": account.Id}
	updateAccount := bson.M{
		"updatedAt": time.Now().Unix(),
	}
	if account.Email != "" {
		updateAccount["email"] = account.Email
	}
	if account.Name != "" {
		updateAccount["name"] = account.Name
	}
	if account.Password != "" {
		updateAccount["password"] = helpers.GeneratePasswordHash([]byte(account.Password))
	}
	update := bson.M{"$set": updateAccount}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Update account successful"})

}

func (accountController AccountController) DeleteAccount(c *gin.Context) {

}

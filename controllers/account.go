package controllers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"ima-svc-management/helpers"
	"ima-svc-management/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountController struct {
	MongoClient *mongo.Client
}

func InitAccount(mongoClient *mongo.Client) *AccountController {
	return &AccountController{
		MongoClient: mongoClient,
	}
}

// @Summary Add account
// @Description create new account
// @Param body body model.AccountModel true "body"
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/account/add [post]
func (accountController AccountController) AddAccount(c *gin.Context) {

	collection := accountController.MongoClient.Database("test").Collection("account")

	account := model.AccountModel{}
	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	dataAccount := bson.M{
		"name":      account.Name,
		"email":     account.Email,
		"password":  helpers.GeneratePasswordHash([]byte(account.Password)),
		"role":      account.Role,
		"createdAt": time.Now().Unix(),
		"updatedAt": nil,
	}

	hashId, err := bson.Marshal(dataAccount)
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.Sum(hashId)

	dataAccount["_id"] = hex.EncodeToString(hash[:])

	registeredAccount := model.AccountModel{}
	filter := bson.M{"email": account.Email}
	err = collection.FindOne(context.TODO(), filter).Decode(&registeredAccount)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
		c.Abort()
		return
	}

	_, err = collection.InsertOne(context.Background(), dataAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Create account successful"})

}

// @Summary Get all account
// @Description get all account with pagination
// @Param body body model.PaginateRoleModel true "body"
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,datas=[]model.AccountModel} "ok"
// @Router /api/v1/account/getAll [post]
// @Security BearerAuth
func (accountController AccountController) GetAccount(c *gin.Context) {

	paginationModel := model.PaginateAccountModel{}

	err := c.BindJSON(&paginationModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := accountController.MongoClient.Database("test").Collection("account")

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
			"role":       account.Role,
			"created_at": account.CreatedAt,
			"updated_at": account.UpdatedAt,
		}
		datas = append(datas, data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})

}

// @Summary Get account by email
// @Description get account using email
// @Param email query string true "email"
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,datas=[]model.AccountModel} "ok"
// @Router /api/v1/account/getByEmail [get]
// @Security BearerAuth
func (accountController AccountController) GetAccountByEmail(c *gin.Context) {

	account := model.AccountModel{}

	email := c.Query("email")

	collection := accountController.MongoClient.Database("test").Collection("account")

	filter := bson.M{"email": email}

	err := collection.FindOne(context.TODO(), filter).Decode(&account)
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
		"role":       account.Role,
		"created_at": account.CreatedAt,
		"updated_at": account.UpdatedAt,
	}
	datas = append(datas, data)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

// @Summary Get account by id
// @Description get account using id
// @Param id query string true "id"
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,datas=[]model.AccountModel} "ok"
// @Router /api/v1/account/getById [get]
// @Security BearerAuth
func (accountController AccountController) GetAccountById(c *gin.Context) {

	account := model.AccountModel{}

	id := c.Query("id")

	collection := accountController.MongoClient.Database("test").Collection("account")

	filter := bson.M{"_id": id}

	err := collection.FindOne(context.TODO(), filter).Decode(&account)
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
		"role":       account.Role,
		"created_at": account.CreatedAt,
		"updated_at": account.UpdatedAt,
	}
	datas = append(datas, data)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

// @Summary Update account
// @Description Update account
// @Param body body model.AccountModel true "body"
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/account/update [put]
// @Security BearerAuth
func (accountController AccountController) UpdateAccount(c *gin.Context) {

	collection := accountController.MongoClient.Database("test").Collection("account")

	account := model.AccountModel{}
	err := c.BindJSON(&account)
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
	if account.Role != "" {
		updateAccount["role"] = account.Role
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

// @Summary Delete account by id
// @Description delete account using id
// @Param id query string true "id"
// @Tags Account
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/account/delete [delete]
// @Security BearerAuth
func (accountController AccountController) DeleteAccount(c *gin.Context) {
	id := c.Query("id")
	collection := accountController.MongoClient.Database("test").Collection("account")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Delete account successful"})

}

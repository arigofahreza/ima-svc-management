package controllers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"ima-svc-management/config"
	"ima-svc-management/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoleController struct{}

func (roleController RoleController) AddRole(c *gin.Context) {
	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("role")

	role := model.RoleModel{}
	err = c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	dataRole := bson.M{
		"name":        role.Name,
		"role":        role.Role,
		"description": role.Description,
		"createdAt":   time.Now().Unix(),
		"updatedAt":   nil,
	}

	hashId, err := bson.Marshal(dataRole)
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.Sum(hashId)

	dataRole["_id"] = hex.EncodeToString(hash[:])

	_, err = collection.InsertOne(context.Background(), dataRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Create role successful"})
}

func (roleController RoleController) GetRole(c *gin.Context) {
	paginationModel := model.PaginateRoleModel{}

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
	collection := mongoClient.Database("test").Collection("role")

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
		role := model.RoleModel{}
		if err := cursor.Decode(&role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		data := map[string]interface{}{
			"id":          role.Id,
			"name":        role.Name,
			"role":        role.Role,
			"description": role.Description,
			"created_at":  role.CreatedAt,
			"updated_at":  role.UpdatedAt,
		}
		datas = append(datas, data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

func (roleController RoleController) GetRoleById(c *gin.Context) {
	role := model.RoleModel{}

	id := c.Query("id")

	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("role")

	filter := bson.M{"_id": id}

	err = collection.FindOne(context.TODO(), filter).Decode(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	datas := make([]map[string]interface{}, 0)
	data := map[string]interface{}{
		"id":          role.Id,
		"name":        role.Name,
		"role":        role.Role,
		"description": role.Description,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
	}
	datas = append(datas, data)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

func (roleController RoleController) UpdateRole(c *gin.Context) {

	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("role")

	role := model.RoleModel{}
	err = c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	filter := bson.M{"_id": role.Id}
	updateRole := bson.M{
		"updatedAt": time.Now().Unix(),
	}
	if role.Description != "" {
		updateRole["description"] = role.Description
	}
	if role.Name != "" {
		updateRole["name"] = role.Name
	}
	if role.Role != ""{
		updateRole["role"] = role.Role
	}
	update := bson.M{"$set": updateRole}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Update role successful"})

}

func (roleController RoleController) DeleteRole(c *gin.Context) {
	id := c.Query("id")
	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	collection := mongoClient.Database("test").Collection("role")

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Delete role successful"})
}

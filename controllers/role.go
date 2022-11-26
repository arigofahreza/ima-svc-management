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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoleController struct {
	MongoClient *mongo.Client
}

func InitRole() *RoleController {
	mongoClient, err := config.InitMongo().Mongo()
	if err != nil {
		return nil
	}
	return &RoleController{
		MongoClient: mongoClient,
	}
}

// @Summary Add role
// @Description create new role
// @Param body body model.RoleModel true "body"
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/role/add [post]
// @Security ApiKeyAuth
func (roleController RoleController) AddRole(c *gin.Context) {
	collection := roleController.MongoClient.Database("test").Collection("role")

	role := model.RoleModel{}
	err := c.BindJSON(&role)
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

// @Summary Get all role
// @Description get all role with pagination
// @Param body body model.PaginateRoleModel true "body"
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,datas=[]model.RoleModel} "ok"
// @Router /api/v1/role/getAll [post]
// @Security ApiKeyAuth
func (roleController RoleController) GetRole(c *gin.Context) {
	paginationModel := model.PaginateRoleModel{}

	err := c.BindJSON(&paginationModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	collection := roleController.MongoClient.Database("test").Collection("role")

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
			"createdAt":   role.CreatedAt,
			"updatedAt":   role.UpdatedAt,
		}
		datas = append(datas, data)
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

// @Summary Get role by id
// @Description get role using id
// @Param id query string true "id"
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,datas=[]model.RoleModel} "ok"
// @Router /api/v1/role/getById [get]
// @Security ApiKeyAuth
func (roleController RoleController) GetRoleById(c *gin.Context) {
	role := model.RoleModel{}

	id := c.Query("id")

	collection := roleController.MongoClient.Database("test").Collection("role")

	filter := bson.M{"_id": id}

	err := collection.FindOne(context.TODO(), filter).Decode(&role)
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
		"createdAt":   role.CreatedAt,
		"updatedAt":   role.UpdatedAt,
	}
	datas = append(datas, data)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "data": datas})
}

// @Summary Update role
// @Description Update role
// @Param body body model.RoleModel true "body"
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/role/update [put]
// @Security ApiKeyAuth
func (roleController RoleController) UpdateRole(c *gin.Context) {

	collection := roleController.MongoClient.Database("test").Collection("role")

	role := model.RoleModel{}
	err := c.BindJSON(&role)
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
	if role.Role != "" {
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

// @Summary Delete role by id
// @Description delete role using id
// @Param id query string true "id"
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} object{status=string,message=string} "ok"
// @Router /api/v1/role/delete [delete]
// @Security ApiKeyAuth
func (roleController RoleController) DeleteRole(c *gin.Context) {
	id := c.Query("id")

	collection := roleController.MongoClient.Database("test").Collection("role")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Delete role successful"})
}

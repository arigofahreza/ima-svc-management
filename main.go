package main

import (
	"context"
	"ima-svc-management/config"
	"ima-svc-management/controllers"
	docs "ima-svc-management/docs"
	"ima-svc-management/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title IMA Reprocess Management API
// @version 1.0
// @description API for management account and role IMA Reprocess Project
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization

func main() {

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())
	docs.SwaggerInfo.BasePath = "/"
	mongoClient, err := config.Mongo()
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.TODO())
	redisClient, err := config.Redis()
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()
	accountController := controllers.InitAccount(config.MongoClient)
	roleController := controllers.InitRole(config.MongoClient)
	authController := controllers.InitAuth(config.RedisClient, config.MongoClient)

	mainGroup := router.Group("/api/v1")
	{
		account := mainGroup.Group("/account")
		{
			account.POST("/add", accountController.AddAccount)
			account.GET("/getById", AuthMiddleware(), accountController.GetAccountById)
			account.GET("/getByEmail", AuthMiddleware(), accountController.GetAccountByEmail)
			account.POST("/getAll", AuthMiddleware(), accountController.GetAccount)
			account.PUT("/update", AuthMiddleware(), accountController.UpdateAccount)
			account.DELETE("/delete", AuthMiddleware(), accountController.DeleteAccount)
		}

		role := mainGroup.Group("/role")
		{
			role.POST("/add", AuthMiddleware(), roleController.AddRole)
			role.GET("/getById", AuthMiddleware(), roleController.GetRoleById)
			role.POST("/getAll", AuthMiddleware(), roleController.GetRole)
			role.PUT("/update", AuthMiddleware(), roleController.UpdateRole)
			role.DELETE("/delete", AuthMiddleware(), roleController.DeleteRole)
		}

		auth := mainGroup.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout", AuthMiddleware(), authController.Logout)
			auth.GET("/refresh", authController.Refresh)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(":45541")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := helpers.Auth{}
		err := auth.TokenValid(c)
		if err != nil && err.Error() == "Token is expired" {
			c.JSON(http.StatusUnauthorized, "Token is expired")
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusUnauthorized, "Invalid Token")
			c.Abort()
			return
		}
		c.Next()
	}
}

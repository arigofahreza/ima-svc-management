package main

import (
	"ima-svc-management/config"
	"ima-svc-management/controllers"
	docs "ima-svc-management/docs"

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
	err := config.Mongo()
	if err != nil {
		panic(err)
	}
	err = config.Redis()
	if err != nil {
		panic(err)
	}
	accountController := controllers.InitAccount(config.MongoClient)
	roleController := controllers.InitRole(config.MongoClient)
	authController := controllers.InitAuth(config.RedisClient, config.MongoClient)

	mainGroup := router.Group("/api/v1")
	{
		account := mainGroup.Group("/account")
		{
			account.POST("/add", accountController.AddAccount)
			account.GET("/getById", accountController.GetAccountById)
			account.POST("/getAll", accountController.GetAccount)
			account.PUT("/update", accountController.UpdateAccount)
			account.DELETE("/delete", accountController.DeleteAccount)
			account.GET("/checkEmail", accountController.CheckEmail)
		}

		role := mainGroup.Group("/role")
		{
			role.POST("/add", roleController.AddRole)
			role.GET("/getById", roleController.GetRoleById)
			role.POST("/getAll", roleController.GetRole)
			role.PUT("/update", roleController.UpdateRole)
			role.DELETE("/delete", roleController.DeleteRole)
		}

		auth := mainGroup.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authController.Logout)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(":8000")

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

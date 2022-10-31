package main

import (
	"ima-svc-management/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())

	accountController := controllers.AccountController{}

	mainGroup := router.Group("/api/v1")
	{
		account := mainGroup.Group("/account")
		{
			account.POST("/add", accountController.AddAccount)
			account.GET("/getById", accountController.GetAccountById)
			account.POST("/getAll", accountController.GetAccount)
			account.PUT("/update", accountController.UpdateAccount)
		}
	}

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

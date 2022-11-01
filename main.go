package main

import (
	"ima-svc-management/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

func main() {

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())

	accountController := controllers.AccountController{}
	roleController := controllers.RoleController{}

	mainGroup := router.Group("/api/v1")
	{
		account := mainGroup.Group("/account")
		{
			account.POST("/add", accountController.AddAccount)
			account.GET("/getById", accountController.GetAccountById)
			account.POST("/getAll", accountController.GetAccount)
			account.PUT("/update", accountController.UpdateAccount)
			account.DELETE("/delete", accountController.DeleteAccount)
		}

		role := mainGroup.Group("/role")
		{
			role.POST("/add", roleController.AddRole)
			role.GET("/getById", roleController.GetRoleById)
			role.POST("/getAll", roleController.GetRole)
			role.PUT("/update", roleController.UpdateRole)
			role.DELETE("/delete", roleController.DeleteRole)
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

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userkey")
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

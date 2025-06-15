package routes

import (
	"ecommerce/controllers"
	"ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	user := r.Group("/api")
	{
		user.POST("/register", controllers.Register)
		user.POST("/login", controllers.Login)
		user.GET("/profile", middleware.AuthMiddleware(), controllers.Profile)
	}
}

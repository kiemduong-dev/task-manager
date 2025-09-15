package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/kiemduong-dev/task-manager/controllers"
	"github.com/kiemduong-dev/task-manager/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// protected
	auth := r.Group("/", middleware.AuthMiddleware(), middleware.RateLimit())
	{
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.ListTasks)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	category := r.Group("/categories")
	{
		category.POST("/", controllers.CreateCategory)
		category.GET("/", controllers.GetCategories)
		category.GET("/:id", controllers.GetCategory)
		category.PUT("/:id", controllers.UpdateCategory)
		category.DELETE("/:id", controllers.DeleteCategory)
	}
	return r
}

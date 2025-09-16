package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/kiemduong-dev/task-manager/controllers"
	"github.com/kiemduong-dev/task-manager/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID(), middleware.RequestLogger())

	r.POST("/register", controllers.Register)
	// apply rate limit to login
	r.POST("/login", middleware.RateLimit(), controllers.Login)

	// protected
	auth := r.Group("/", middleware.AuthMiddleware(), middleware.RateLimit())
	{
		auth.POST("/tasks", controllers.CreateTask)
		auth.GET("/tasks", controllers.ListTasks)
		auth.GET("/tasks/:id", controllers.GetTask)
		auth.PATCH("/tasks/:id/complete", controllers.CompleteTask)
		auth.PUT("/tasks/:id", controllers.UpdateTask)
		auth.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	category := r.Group("/categories", middleware.AuthMiddleware())
	{
		category.GET("/", controllers.GetCategories)
		category.GET("/:id", controllers.GetCategory)
		adminCat := category.Group("/", middleware.Authorize("admin"))
		{
			adminCat.POST("/", controllers.CreateCategory)
			adminCat.PUT(":id", controllers.UpdateCategory)
			adminCat.DELETE(":id", controllers.DeleteCategory)
		}
	}

	// serve static OpenAPI at /openapi.json if file present
	r.StaticFile("/openapi.json", "openapi.json")
	return r
}

package routes

import (
	"TaskManager/handlers"
	"TaskManager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Public authentication routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)
		auth.POST("/refresh", handlers.RefreshToken)
	}

	// Protected task routes
	tasks := api.Group("/tasks")
	tasks.Use(middleware.RequireAuth)
	{
		tasks.GET("", handlers.GetTasks)
		tasks.POST("", handlers.CreateTask)
		tasks.GET("/:id", handlers.GetTask)
		tasks.PUT("/:id", handlers.UpdateTask)
		tasks.PATCH("/:id", handlers.PatchTask)
		tasks.DELETE("/:id", handlers.DeleteTask)
	}
}

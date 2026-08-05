package main

import (
	"TaskManager/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/tasks", handlers.GetTasks)
		api.POST("/tasks", handlers.CreateTask)
		api.GET("/tasks/:id", handlers.GetTask)
	}

	r.Run(":8080")
}

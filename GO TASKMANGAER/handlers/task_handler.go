package handlers

import (
	"net/http"
	"strconv"

	"TaskManager/db"
	"TaskManager/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	var tasks []models.Task
	query := db.DB.Model(&models.Task{}).Where("user_id = ?", user.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Pagination support (?page=1&limit=10)
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			limit := 10
			if limitStr := c.Query("limit"); limitStr != "" {
				if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
					limit = l
				}
			}
			offset := (page - 1) * limit
			query = query.Offset(offset).Limit(limit)
		}
	} else if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			query = query.Limit(limit)
		}
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func CreateTask(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	var input models.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		DueDate:     input.DueDate,
		UserID:      user.ID,
	}

	if err := db.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func GetTask(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	var task models.Task
	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func UpdateTask(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	var task models.Task
	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != "" && !input.Status.IsValid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid status: must be pending, in_progress, or done"})
		return
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	task.Description = input.Description
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}

	if err := db.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func PatchTask(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	var task models.Task
	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	var input models.PatchTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != nil && !input.Status.IsValid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid status: must be pending, in_progress, or done"})
		return
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		task.Status = *input.Status
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}

	if err := db.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func DeleteTask(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user := currentUser.(models.User)

	var task models.Task
	if err := db.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	if err := db.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

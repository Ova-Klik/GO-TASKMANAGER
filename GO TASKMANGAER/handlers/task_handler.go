package handlers


import (
    "net/http"
    "strconv"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "TaskManager/models"
)

var (
    tasks  []models.Task
    nextID uint = 1
    mu     sync.Mutex
)

func GetTasks(c *gin.Context) {
    mu.Lock()
    defer mu.Unlock()
    c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func CreateTask(c *gin.Context) {
    var input models.CreateTaskInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    mu.Lock()
    newTask := models.Task{
        ID:          nextID,
        Title:       input.Title,
        Description: input.Description,
        Status:      models.StatusPending,
        CreatedAt:   time.Now(),
    }
    nextID++
    tasks = append(tasks, newTask)
    mu.Unlock()

    c.JSON(http.StatusCreated, gin.H{"data": newTask})
}

func GetTask(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    mu.Lock()
    defer mu.Unlock()
    for _, t := range tasks {
        if t.ID == uint(id) {
            c.JSON(http.StatusOK, gin.H{"data": t})
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

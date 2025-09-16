package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kiemduong-dev/task-manager/config"
	"github.com/kiemduong-dev/task-manager/models"
)

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"omitempty"`
	CategoryID  uint       `json:"category_id" binding:"required,gt=0"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title" binding:"omitempty,min=1"`
	Description *string    `json:"description" binding:"omitempty"`
	CategoryID  *uint      `json:"category_id" binding:"omitempty,gt=0"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
	Completed   *bool      `json:"completed" binding:"omitempty"`
}

// CreateTask
func CreateTask(c *gin.Context) {
	var in CreateTaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// validate category exists
	var cat models.Category
	if err := config.DB.First(&cat, in.CategoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	uid := c.GetUint("user_id")
	task := models.Task{
		Title:       in.Title,
		Description: in.Description,
		CategoryID:  in.CategoryID,
		DueDate:     time.Time{},
		Completed:   false,
		UserID:      uid,
	}
	if in.DueDate != nil {
		task.DueDate = *in.DueDate
	}
	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create task failed"})
		return
	}
	log.Printf("user_id=%d action=create_task id=%d title=%q", uid, task.ID, task.Title)
	c.JSON(http.StatusCreated, task)
}

// ListTasks (own user's tasks) or all for admin
func ListTasks(c *gin.Context) {
	role := c.GetString("role")
	// filters
	completed := c.Query("completed") // "true" | "false" | ""
	categoryID := c.Query("category_id")
	dueFrom := c.Query("due_from") // RFC3339
	dueTo := c.Query("due_to")
	// pagination
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageStr := c.DefaultQuery("page", "1")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	dbq := config.DB.Model(&models.Task{}).Preload("Category")
	if role != "admin" {
		uid := c.GetUint("user_id")
		dbq = dbq.Where("user_id = ?", uid)
	}
	if completed == "true" {
		dbq = dbq.Where("completed = ?", true)
	} else if completed == "false" {
		dbq = dbq.Where("completed = ?", false)
	}
	if categoryID != "" {
		dbq = dbq.Where("category_id = ?", categoryID)
	}
	if dueFrom != "" {
		if t, err := time.Parse(time.RFC3339, dueFrom); err == nil {
			dbq = dbq.Where("due_date >= ?", t)
		}
	}
	if dueTo != "" {
		if t, err := time.Parse(time.RFC3339, dueTo); err == nil {
			dbq = dbq.Where("due_date <= ?", t)
		}
	}
	var total int64
	dbq.Count(&total)
	var tasks []models.Task
	dbq.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// UpdateTask
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	uid := c.GetUint("user_id")
	role := c.GetString("role")
	if task.UserID != uid && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var in UpdateTaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Description != nil {
		task.Description = *in.Description
	}
	if in.CategoryID != nil {
		var cat models.Category
		if err := config.DB.First(&cat, *in.CategoryID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		task.CategoryID = *in.CategoryID
	}
	if in.DueDate != nil {
		task.DueDate = *in.DueDate
	}
	if in.Completed != nil {
		task.Completed = *in.Completed
	}
	task.UpdatedAt = time.Now()
	config.DB.Save(&task)
	log.Printf("user_id=%d action=update_task id=%d", uid, task.ID)
	c.JSON(http.StatusOK, task)
}

// DeleteTask
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	uid := c.GetUint("user_id")
	role := c.GetString("role")
	if task.UserID != uid && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	config.DB.Delete(&task)
	log.Printf("user_id=%d action=delete_task id=%d", uid, task.ID)
	c.JSON(http.StatusNoContent, nil)
}

// GetTask returns a task by id with ownership/admin check
func GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var task models.Task
	if err := config.DB.Preload("Category").First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	uid := c.GetUint("user_id")
	role := c.GetString("role")
	if task.UserID != uid && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CompleteTask sets completed state
func CompleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	uid := c.GetUint("user_id")
	role := c.GetString("role")
	if task.UserID != uid && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var body struct {
		Completed bool `json:"completed" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.Completed = body.Completed
	task.UpdatedAt = time.Now()
	config.DB.Save(&task)
	log.Printf("user_id=%d action=complete_task id=%d completed=%v", uid, task.ID, task.Completed)
	c.JSON(http.StatusOK, task)
}

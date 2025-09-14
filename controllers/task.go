package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kiemduong-dev/task-manager/config"
	"github.com/kiemduong-dev/task-manager/models"
)

// CreateTask
func CreateTask(c *gin.Context) {
	var in models.Task
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := c.GetUint("user_id")
	in.UserID = uid
	if err := config.DB.Create(&in).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create task failed"})
		return
	}
	c.JSON(http.StatusCreated, in)
}

// ListTasks (own user's tasks) or all for admin
func ListTasks(c *gin.Context) {
	role := c.GetString("role")
	var tasks []models.Task
	if role == "admin" {
		config.DB.Find(&tasks)
	} else {
		uid := c.GetUint("user_id")
		config.DB.Where("user_id = ?", uid).Find(&tasks)
	}
	c.JSON(http.StatusOK, tasks)
}

// UpdateTask
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error":"not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"})
		return
	}
	uid := c.GetUint("user_id")
	role := c.GetString("role")
	if task.UserID != uid && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error":"forbidden"})
		return
	}
	var in models.Task
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	task.Title = in.Title
	task.Description = in.Description
	task.Category = in.Category
	task.DueDate = in.DueDate
	task.Completed = in.Completed
	task.UpdatedAt = time.Now()
	config.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

// DeleteTask
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error":"not found"})
		return
	}
	uid := c.GetUint("user_id")
	role := c.GetString("role")
	if task.UserID != uid && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error":"forbidden"})
		return
	}
	config.DB.Delete(&task)
	c.JSON(http.StatusNoContent, nil)
}

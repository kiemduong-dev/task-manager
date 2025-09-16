package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiemduong-dev/task-manager/config"
	"github.com/kiemduong-dev/task-manager/models"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required,min=1"`
}

type UpdateCategoryInput struct {
	Name *string `json:"name" binding:"omitempty,min=1"`
}

func CreateCategory(c *gin.Context) {
	var in CreateCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category := models.Category{Name: in.Name}
	if err := config.DB.Create(&category).Error; err != nil {
		// Trả 409 nếu vi phạm unique (duplicate name)
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "category name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("user_id=%d action=create_category id=%d name=%s", c.GetUint("user_id"), category.ID, category.Name)
	c.JSON(http.StatusCreated, category)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	var category models.Category
	id := c.Param("id")
	if err := config.DB.Preload("Tasks").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	var in UpdateCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if in.Name != nil {
		category.Name = *in.Name
	}
	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("user_id=%d action=update_category id=%d name=%s", c.GetUint("user_id"), category.ID, category.Name)
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("user_id=%d action=delete_category id=%s", c.GetUint("user_id"), id)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

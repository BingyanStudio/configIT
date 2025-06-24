package admin

import (
	"net/http"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/gin-gonic/gin"
)

// GetApps retrieves all applications
func GetApps(c *gin.Context) {
	var apps []model.App

	if err := model.DB().Preload("Namespace").Preload("Owner").Preload("Department").Preload("Scopes").Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	c.JSON(http.StatusOK, apps)
}

// GetApp retrieves a specific application by name
func GetApp(c *gin.Context) {
	name := c.Param("name")

	app, err := model.GetAppByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Load related data
	if err := model.DB().Preload("Namespace").Preload("Owner").Preload("Department").Preload("Scopes").First(app, app.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve application details"})
		return
	}

	c.JSON(http.StatusOK, app)
}

// CreateApp creates a new application
func CreateApp(c *gin.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate app data
	if app.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application name is required"})
		return
	}

	// Set default values if not provided
	if app.EditPermission == "" {
		app.EditPermission = model.EditPermisionOwner
	}

	if app.AutoRefresh == "" {
		app.AutoRefresh = model.AutoRefreshPassive
	}

	if err := model.InsertApp(c.Request.Context(), app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, app)
}

// UpdateApp updates an existing application
func UpdateApp(c *gin.Context) {
	name := c.Param("name")

	// Check if app exists
	existingApp, err := model.GetAppByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	var updateData model.App
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve the ID and Name
	updateData.ID = existingApp.ID
	updateData.Name = existingApp.Name

	if err := model.UpdateApp(c.Request.Context(), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, updateData)
}

// DeleteApp deletes an application
func DeleteApp(c *gin.Context) {
	name := c.Param("name")

	// Check if app exists
	app, err := model.GetAppByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if err := model.DeleteApp(c.Request.Context(), *app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}

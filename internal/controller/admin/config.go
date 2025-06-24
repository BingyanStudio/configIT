package admin

import (
	"net/http"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/gin-gonic/gin"
)

// GetConfigs retrieves all configurations for an app
func GetConfigs(c *gin.Context) {
	appName := c.Param("app_name")

	var configs []model.Config
	if err := model.DB().Where("app = ?", appName).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configurations"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetConfig retrieves a specific configuration by key
func GetConfig(c *gin.Context) {
	appName := c.Param("app_name")
	configKey := c.Param("config_key")

	config, err := model.GetConfig(c.Request.Context(), appName, configKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// CreateConfig creates a new configuration
func CreateConfig(c *gin.Context) {
	appName := c.Param("app_name")

	// Check if app exists
	app, err := model.GetAppByName(c.Request.Context(), appName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	var config model.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the app for the config
	config.App = *app

	// Check if a config with the same key already exists
	var existingConfig model.Config
	if err := model.DB().Where("app = ? AND key = ?", appName, config.Key).First(&existingConfig).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Configuration with this key already exists"})
		return
	}

	if err := model.DB().Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// UpdateConfig updates an existing configuration
func UpdateConfig(c *gin.Context) {
	appName := c.Param("app_name")
	configKey := c.Param("config_key")

	// Check if config exists
	existingConfig, err := model.GetConfig(c.Request.Context(), appName, configKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	var updateData model.Config
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields but keep the ID and App
	existingConfig.Type = updateData.Type
	existingConfig.Value = updateData.Value
	existingConfig.Comment = updateData.Comment

	if err := model.UpdateConfig(c.Request.Context(), existingConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, existingConfig)
}

// DeleteConfig deletes a configuration
func DeleteConfig(c *gin.Context) {
	appName := c.Param("app_name")
	configKey := c.Param("config_key")

	// Check if config exists
	existingConfig, err := model.GetConfig(c.Request.Context(), appName, configKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	if err := model.DB().Delete(existingConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}

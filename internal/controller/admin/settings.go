package admin

import (
	"net/http"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/gin-gonic/gin"
)

// GetAllSettings retrieves all settings
func GetAllSettings(c *gin.Context) {
	var settings []model.Settings

	if err := model.DB().Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve settings"})
		return
	}

	// Convert to a map for easier consumption by clients
	settingsMap := make(map[string]string)
	for _, setting := range settings {
		// Don't expose sensitive settings like JWTSecret
		if setting.Key == "JWTSecret" {
			continue
		}
		settingsMap[setting.Key] = setting.Value
	}

	c.JSON(http.StatusOK, settingsMap)
}

// GetSetting retrieves a specific setting by key
func GetSetting(c *gin.Context) {
	key := c.Param("key")

	// Don't allow accessing sensitive settings directly
	if key == "JWTSecret" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot access this setting directly"})
		return
	}

	value, err := model.GetSettings(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{key: value})
}

// UpdateSetting updates a specific setting
func UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var data struct {
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if setting exists
	_, err := model.GetSettings(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	if err := model.SetSettings(c.Request.Context(), key, data.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully"})
}

// UpdateSettings updates multiple settings at once
func UpdateSettings(c *gin.Context) {
	var settings map[string]string

	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// Update each setting
	errors := make(map[string]string)
	for key, value := range settings {
		if err := model.SetSettings(ctx, key, value); err != nil {
			errors[key] = err.Error()
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusMultiStatus, gin.H{
			"message": "Some settings could not be updated",
			"errors":  errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All settings updated successfully"})
}

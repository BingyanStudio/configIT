package admin

import (
	"net/http"
	"strconv"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/gin-gonic/gin"
)

// GetAccessScopes retrieves all access scopes
func GetAccessScopes(c *gin.Context) {
	var accessScopes []model.AccessScope

	if err := model.DB().Find(&accessScopes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve access scopes"})
		return
	}

	c.JSON(http.StatusOK, accessScopes)
}

// GetAccessScope retrieves a specific access scope by ID
func GetAccessScope(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access scope ID"})
		return
	}

	var accessScope model.AccessScope
	if err := model.DB().First(&accessScope, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access scope not found"})
		return
	}

	c.JSON(http.StatusOK, accessScope)
}

// CreateAccessScope creates a new access scope
func CreateAccessScope(c *gin.Context) {
	var accessScope model.AccessScope
	if err := c.ShouldBindJSON(&accessScope); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := model.InsertAccessScope(c.Request.Context(), &accessScope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access scope"})
		return
	}

	accessScope.ID = id
	c.JSON(http.StatusCreated, accessScope)
}

// UpdateAccessScope updates an existing access scope
func UpdateAccessScope(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access scope ID"})
		return
	}

	var accessScope model.AccessScope
	if err := model.DB().First(&accessScope, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access scope not found"})
		return
	}

	if err := c.ShouldBindJSON(&accessScope); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.UpdateAccessScope(c.Request.Context(), &accessScope); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update access scope"})
		return
	}

	c.JSON(http.StatusOK, accessScope)
}

// DeleteAccessScope deletes an access scope
func DeleteAccessScope(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access scope ID"})
		return
	}

	var accessScope model.AccessScope
	if err := model.DB().First(&accessScope, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Access scope not found"})
		return
	}

	if err := model.DeleteAccessScope(c.Request.Context(), &accessScope); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete access scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access scope deleted successfully"})
}

package admin

import (
	"net/http"
	"strconv"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/BingyanStudio/configIT/internal/utils"
	"github.com/gin-gonic/gin"
)

// GetNamespaces retrieves all namespaces
func GetNamespaces(c *gin.Context) {

	ns, err := model.GetNamespaces(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve namespaces"})
		return
	}

	c.JSON(http.StatusOK, ns)
}

// GetNamespace retrieves a specific namespace by ID
func GetNamespace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid namespace ID"})
		return
	}

	n, err := model.GetNamespace(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Namespace not found"})
		return
	}

	c.JSON(http.StatusOK, n)
}

// CreateNamespace creates a new namespace
func CreateNamespace(c *gin.Context) {
	var namespace model.Namespace
	if err := c.ShouldBindJSON(&namespace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.InsertNamespace(c, namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create namespace"})
		return
	}

	c.JSON(http.StatusCreated, namespace)
}

// UpdateNamespace updates an existing namespace
func UpdateNamespace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid namespace ID"})
		return
	}

	var namespace model.Namespace
	if _, err := model.GetNamespace(c, uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Namespace not found"})
		return
	}

	if err := c.ShouldBindJSON(&namespace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace.ID = uint(id)

	if err := model.UpdateNamespace(c, namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update namespace"})
		return
	}

	c.JSON(http.StatusOK, namespace)
}

// DeleteNamespace deletes a namespace
func DeleteNamespace(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid namespace ID"})
		return
	}

	n, err := model.GetNamespace(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Namespace not found"})
		return
	}

	if err := model.DeleteNamespace(c, n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete namespace"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Namespace deleted successfully"})
}

// SyncNamespaces synchronizes namespaces from Kubernetes
func SyncNamespaces(c *gin.Context) {
	namespaces, err := utils.GetNamespaces(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get namespaces from Kubernetes"})
		return
	}

	if err := model.ImportNamespaceFromK8s(c, namespaces); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import namespaces from Kubernetes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Namespaces synchronized successfully"})
}

package client

import (
	"net/http"
	"time"

	"github.com/BingyanStudio/configIT/internal/controller"
	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/gin-gonic/gin"
)

// GetAppInfo retrieves and returns application information in JSON format
func GetAppInfo(c *gin.Context) {
	ctx := c.Request.Context()

	// Get app name from request parameters
	appName := c.Param("name")
	if appName == "" {
		controller.Error(c, http.StatusBadRequest, "App name is required")
		return
	}

	// Retrieve app information from database
	app, err := model.GetAppByName(ctx, appName)
	if err != nil {
		controller.ErrNotFoundOrInternal(c, err)
		return
	}

	type AppResponse struct {
		ID               uint   `json:"id"`
		Name             string `json:"name"`
		Namespace        string `json:"namespace"`
		AutoRefresh      string `json:"autoRefresh"`
		AutoRefreshParam string `json:"autoRefreshParam,omitempty"`
		CreatedAt        string `json:"createdAt"`
		UpdatedAt        string `json:"updatedAt"`
	}

	// Map model data to response
	response := AppResponse{
		ID:               app.ID,
		Name:             app.Name,
		Namespace:        app.Namespace.Name,
		AutoRefresh:      app.AutoRefresh,
		AutoRefreshParam: app.AutoRefreshParam,
		CreatedAt:        app.CreatedAt.Local().Format(time.RFC3339),
		UpdatedAt:        app.UpdatedAt.Local().Format(time.RFC3339),
	}

	// Return response
	controller.OK(c, response)
}

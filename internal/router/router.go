package router

import (
	"github.com/BingyanStudio/configIT/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r gin.IRouter) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Type"},
		AllowCredentials: true,
	},
	))

	apiGroup := r.Group("/api/:app_name")
	clientGroup.Use(middleware.Auth)
	clientGroup.GET("", controller.ClientGetConfigs)
	clientGroup.GET("/:key", controller.ClientGetConfig)
	clientGroup.GET("/:key/raw", controller.ClientGetConfigRaw)

	r.POST("/admin/login", controller.Login)
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminAuth)
	adminGroup.GET("/current-user", controller.GetCurrentUser)

	applicationGroup := adminGroup.Group("/applications")
	applicationGroup.GET("", controller.GetApplications)
	applicationGroup.GET("/:app_name", controller.GetApplication)
	applicationGroup.POST("", controller.CreateApplication)
	applicationGroup.PUT("/:app_name", controller.UpdateApplication)
	applicationGroup.DELETE("/:app_name", controller.DeleteApplication)

	versionGroup := adminGroup.Group("/applications/:app_name/versions")
	versionGroup.GET("", controller.GetVersions)
	versionGroup.GET("/:version_name", controller.GetVersion)
	versionGroup.POST("", controller.CreateVersion)
	versionGroup.PUT("/:version_name", controller.UpdateVersion)
	versionGroup.DELETE("/:version_name", controller.DeleteVersion)

	configGroup := adminGroup.Group("/applications/:app_name/versions/:version_name/configs")
	configGroup.GET("", controller.GetConfigs)
	configGroup.GET("/:config_key", controller.GetConfig)
	configGroup.POST("", controller.CreateConfig)
	configGroup.PUT("/:config_key", controller.UpdateConfig)
	configGroup.DELETE("/:config_key", controller.DeleteConfig)
}

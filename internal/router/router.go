package router

import (
	adminController "github.com/BingyanStudio/configIT/internal/controller/admin"
	clientController "github.com/BingyanStudio/configIT/internal/controller/client"
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

	// Client API routes
	clientGroup := r.Group("/api/:app_name")
	clientGroup.Use(middleware.Auth)
	clientGroup.GET("", clientController.GetConfigs)
	clientGroup.GET("/:key", clientController.GetConfig)
	clientGroup.GET("/:key/raw", clientController.GetConfigRaw)

	// Authentication routes
	r.POST("/admin/login", adminController.Login)

	// Admin API routes
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminAuth)
	adminGroup.GET("/current-user", adminController.GetCurrentUser)

	// Application routes
	applicationGroup := adminGroup.Group("/applications")
	applicationGroup.GET("", adminController.GetApps)
	applicationGroup.GET("/:app_name", adminController.GetApp)
	applicationGroup.POST("", adminController.CreateApp)
	applicationGroup.PUT("/:app_name", adminController.UpdateApp)
	applicationGroup.DELETE("/:app_name", adminController.DeleteApp)

	// Config routes
	configGroup := adminGroup.Group("/applications/:app_name/configs")
	configGroup.GET("", adminController.GetConfigs)
	configGroup.GET("/:config_key", adminController.GetConfig)
	configGroup.POST("", adminController.CreateConfig)
	configGroup.PUT("/:config_key", adminController.UpdateConfig)
	configGroup.DELETE("/:config_key", adminController.DeleteConfig)

	// Department routes
	departmentGroup := adminGroup.Group("/departments")
	departmentGroup.GET("", adminController.GetDepartments)
	departmentGroup.GET("/:id", adminController.GetDepartment)
	departmentGroup.GET("/:id/members", adminController.GetDepartmentMembers)
	departmentGroup.POST("", adminController.CreateDepartment)
	departmentGroup.PUT("/:id", adminController.UpdateDepartment)
	departmentGroup.DELETE("/:id", adminController.DeleteDepartment)

	// User routes
	userGroup := adminGroup.Group("/users")
	userGroup.GET("", adminController.GetUsers)
	userGroup.GET("/:id", adminController.GetUser)
	userGroup.GET("/sub/:sub", adminController.GetUserBySub)
	userGroup.POST("", adminController.CreateUser)
	userGroup.PUT("/:id", adminController.UpdateUser)
	userGroup.DELETE("/:id", adminController.DeleteUser)

	// Namespace routes
	namespaceGroup := adminGroup.Group("/namespaces")
	namespaceGroup.GET("", adminController.GetNamespaces)
	namespaceGroup.GET("/:id", adminController.GetNamespace)
	namespaceGroup.POST("", adminController.CreateNamespace)
	namespaceGroup.PUT("/:id", adminController.UpdateNamespace)
	namespaceGroup.DELETE("/:id", adminController.DeleteNamespace)
	namespaceGroup.POST("/sync", adminController.SyncNamespaces)

	// Access scope routes
	accessScopeGroup := adminGroup.Group("/access-scopes")
	accessScopeGroup.GET("", adminController.GetAccessScopes)
	accessScopeGroup.GET("/:id", adminController.GetAccessScope)
	accessScopeGroup.POST("", adminController.CreateAccessScope)
	accessScopeGroup.PUT("/:id", adminController.UpdateAccessScope)
	accessScopeGroup.DELETE("/:id", adminController.DeleteAccessScope)

	// Settings routes
	settingsGroup := adminGroup.Group("/settings")
	settingsGroup.GET("", adminController.GetAllSettings)
	settingsGroup.GET("/:key", adminController.GetSetting)
	settingsGroup.PUT("/:key", adminController.UpdateSetting)
	settingsGroup.PUT("", adminController.UpdateSettings)
}

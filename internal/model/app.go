package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	EditPermissionAnyone     = "anyone"
	EditPermisionOwner       = "owner"
	EditPermissionDepartment = "department"
)

const (
	AutoRefreshConfigMap = "configmap"
	AutoRefreshHook      = "hook"
	AutoRefreshPassive   = "passive"
)

type App struct {
	gorm.Model

	Name             string
	Namespace        Namespace
	Owner            User
	EditPermission   string `gorm:"type:enum('anyone', 'owner', 'department')"`
	Department       Department
	Scope            []AccessScope
	AutoRefresh      string `gorm:"type:enum('configmap', 'hook', 'passive')"`
	AutoRefreshParam string // configmap name or hook url
	InCluster        bool   `gorm:"default:false"`
	LastRefreshTime  time.Time
}

func InsertApp(app App) error {
	return db.Create(&app).Error
}

func DeleteApp(app App) error {
	// First delete the access scopes associated with the app
	if err := db.Model(&app).Association("Scope").Delete(app.Scope); err != nil {
		return err
	}
	return db.Delete(&app).Error
}

func UpdateApp(app App) error {
	return db.Save(&app).Error
}

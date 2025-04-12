package model

import (
	"context"

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
	Scopes           []AccessScope
	AutoRefresh      string `gorm:"type:enum('configmap', 'hook', 'passive')"`
	AutoRefreshParam string // configmap name or hook url
}

func InsertApp(ctx context.Context, app App) error {
	return db.WithContext(ctx).Create(&app).Error
}

func DeleteApp(ctx context.Context, app App) error {
	// First delete the access scopes associated with the app
	if err := db.WithContext(ctx).Model(&app).Association("Scopes").Delete(app.Scopes); err != nil {
		return err
	}
	return db.WithContext(ctx).Delete(&app).Error
}

func UpdateApp(ctx context.Context, app App) error {
	return db.WithContext(ctx).Save(&app).Error
}

func GetAppByName(ctx context.Context, name string) (*App, error) {
	var app App
	if err := db.WithContext(ctx).Where("name = ?", name).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

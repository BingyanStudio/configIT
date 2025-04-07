package model

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Sub        string `gorm:"uniqueIndex"` // OIDC sub or username
	Password   string // password when not using OIDC
	Department Department
}

func InsertUser(user *User) error {
	return db.Create(user).Error
}

func UpdateUser(user *User) error {
	return db.Save(user).Error
}

func DeleteUser(id uint) error {
	// Check if user owns any apps
	var count int64
	if err := db.Model(&App{}).Where("owner = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete user: user still owns %d app(s)", count)
	}
	return db.Delete(&User{}, id).Error
}

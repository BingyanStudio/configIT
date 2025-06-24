package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Sub        string     `gorm:"uniqueIndex"` // OIDC sub or username
	Password   string     // password when not using OIDC
	Department Department `gorm:"foreignKey:name"`
	Role       string     `gorm:"type:enum('admin', 'user')"` // admin or user
}

func InsertUser(ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

func UpdateUser(ctx context.Context, user *User) error {
	return db.WithContext(ctx).Save(user).Error
}

func DeleteUser(ctx context.Context, id uint) error {
	// Check if user owns any apps
	var count int64
	if err := db.WithContext(ctx).Model(&App{}).Where("owner = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete user: user still owns %d app(s)", count)
	}
	return db.WithContext(ctx).Delete(&User{}, id).Error
}

func GetUserBySub(ctx context.Context, sub string) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Where("sub = ?", sub).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers(ctx context.Context, isAdmin bool, dept string) ([]User, error) {
	var users []User
	if isAdmin {
		if err := db.WithContext(ctx).Find(&users).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.WithContext(ctx).Where("department = ?", dept).Find(&users).Error; err != nil {
			return nil, err
		}
	}
	return users, nil
}

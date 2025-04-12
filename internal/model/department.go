package model

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Department struct {
	gorm.Model

	Name string `gorm:"uniqueIndex"`
}

func (d *Department) Members(ctx context.Context) []User {
	var users []User
	db.WithContext(ctx).Model(d).Association("Users").Find(&users)
	return users
}

func InsertDepartment(ctx context.Context, department *Department) (uint, error) {
	if err := db.WithContext(ctx).Create(department).Error; err != nil {
		return 0, err
	}
	return department.ID, nil
}

func DeleteDepartment(ctx context.Context, id uint) error {
	var count int64
	if err := db.WithContext(ctx).Model(&User{}).Where("department_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete department with ID %d: it still has %d members", id, count)
	}
	if err := db.WithContext(ctx).Delete(&Department{}, id).Error; err != nil {
		return err
	}
	return nil
}

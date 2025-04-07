package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Department struct {
	gorm.Model

	Name string `gorm:"uniqueIndex"`
}

func (d *Department) Members() []User {
	var users []User
	db.Model(d).Association("Users").Find(&users)
	return users
}

func InsertDepartment(department *Department) (uint, error) {
	if err := db.Create(department).Error; err != nil {
		return 0, err
	}
	return department.ID, nil
}

func DeleteDepartment(id uint) error {
	var count int64
	if err := db.Model(&User{}).Where("department_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete department with ID %d: it still has %d members", id, count)
	}
	if err := db.Delete(&Department{}, id).Error; err != nil {
		return err
	}
	return nil
}

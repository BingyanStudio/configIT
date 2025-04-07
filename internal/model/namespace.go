package model

import (
	"gorm.io/gorm"
)

type Namespace struct {
	gorm.Model

	Name      string
	InCluster bool `gorm:"default:false"`
}

func InsertNamespace(namespace Namespace) error {
	return db.Create(&namespace).Error
}

func ImportNamespaceFromK8s(namespaces []string) error {
	var existingNamespaces []Namespace
	if err := db.Where("in_cluster = ?", true).Find(&existingNamespaces).Error; err != nil {
		return err
	}

	nsMap := make(map[string]bool)
	for _, ns := range namespaces {
		nsMap[ns] = true
	}

	for _, existingNs := range existingNamespaces {
		if !nsMap[existingNs.Name] {
			if err := db.Model(&existingNs).Update("in_cluster", false).Error; err != nil {
				return err
			}
		} else {
			delete(nsMap, existingNs.Name)
		}
	}

	for ns := range nsMap {
		newNamespace := Namespace{Name: ns, InCluster: true}
		if err := db.Create(&newNamespace).Error; err != nil {
			return err
		}
	}

	return nil
}

func DeleteNamespace(namespace Namespace) error {
	return db.Delete(&namespace).Error
}

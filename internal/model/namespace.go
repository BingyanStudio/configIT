package model

import (
	"context"

	"github.com/BingyanStudio/configIT/internal/utils"
	"gorm.io/gorm"
)

type Namespace struct {
	gorm.Model

	Name        string
	InClusterDB bool `gorm:"default:false"`
}

func InsertNamespace(ctx context.Context, namespace Namespace) error {
	return db.WithContext(ctx).Create(&namespace).Error
}

func ImportNamespaceFromK8s(ctx context.Context, namespaces []string) error {
	var existingNamespaces []Namespace
	if err := db.WithContext(ctx).Where("in_cluster = ?", true).Find(&existingNamespaces).Error; err != nil {
		return err
	}

	nsMap := make(map[string]bool)
	for _, ns := range namespaces {
		nsMap[ns] = true
	}

	for _, existingNs := range existingNamespaces {
		if !nsMap[existingNs.Name] {
			if err := db.WithContext(ctx).Model(&existingNs).Update("in_cluster", false).Error; err != nil {
				return err
			}
		} else {
			delete(nsMap, existingNs.Name)
		}
	}

	for ns := range nsMap {
		newNamespace := Namespace{Name: ns, InClusterDB: true}
		if err := db.WithContext(ctx).Create(&newNamespace).Error; err != nil {
			return err
		}
	}

	return nil
}

func DeleteNamespace(ctx context.Context, namespace Namespace) error {
	return db.WithContext(ctx).Delete(&namespace).Error
}

func UpdateNamespace(ctx context.Context, namespace Namespace) error {
	return db.WithContext(ctx).Save(&namespace).Error
}

func (ns *Namespace) InCluster(ctx context.Context) (bool, error) {
	nss, err := utils.GetNamespaces(ctx)
	if err != nil {
		return false, err
	}
	if !utils.Contains(nss, ns.Name) {
		ns.InClusterDB = false
		return false, db.WithContext(ctx).Save(ns).Error
	}
	return true, nil
}

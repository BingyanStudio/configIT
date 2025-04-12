package model

import (
	"context"

	"gorm.io/gorm"
)

const (
	AccessScopePod       = "pod"
	AccessScopeNamespace = "namespace"
	AccessScopeIP        = "ip"
	AccessScopeToken     = "token"
	AccessScopePublic    = "public"
)

type AccessScope struct {
	gorm.Model

	Scope string `gorm:"type:enum('pod', 'namespace', 'ip', 'token', 'public', 'inherited')"`
	Value string `gorm:"type:text"`
}

func InsertAccessScope(ctx context.Context, scope *AccessScope) (uint, error) {
	if err := db.WithContext(ctx).Create(scope).Error; err != nil {
		return 0, err
	}
	return scope.ID, nil
}

func DeleteAccessScope(ctx context.Context, scope *AccessScope) error {
	return db.WithContext(ctx).Delete(scope).Error
}

func UpdateAccessScope(ctx context.Context, scope *AccessScope) error {
	return db.WithContext(ctx).Save(scope).Error
}

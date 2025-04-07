package model

import (
	"gorm.io/gorm"
)

const (
	AccessScopePod       = "pod"
	AccessScopeNamespace = "namespace"
	AccessScopeIP        = "ip"
	AccessScopeToken     = "token"
	AccessScopePublic    = "public"
	AccessScopeInherited = "inherited"
)

type AccessScope struct {
	gorm.Model

	Scope string `gorm:"type:enum('pod', 'namespace', 'ip', 'token', 'public', 'inherited')"`
	Value string `gorm:"type:text"`
}

func (a *AccessScope) IsInherited() bool {
	return a.Scope == AccessScopeInherited
}

func InsertAccessScope(scope *AccessScope) (uint, error) {
	if err := db.Create(scope).Error; err != nil {
		return 0, err
	}
	return scope.ID, nil
}

func DeleteAccessScope(scope *AccessScope) error {
	return db.Delete(scope).Error
}

func UpdateAccessScope(scope *AccessScope) error {
	return db.Save(scope).Error
}

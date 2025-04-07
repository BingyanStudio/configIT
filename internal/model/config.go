package model

import (
	"gorm.io/gorm"
)

type ConfigType string

const (
	ConfigTypeString ConfigType = "string"
	ConfigTypeNumber ConfigType = "number"
	ConfigTypeBool   ConfigType = "bool"
	ConfigTypeJson   ConfigType = "Json"
)

type Config struct {
	gorm.Model

	App     App
	Type    ConfigType `gorm:"type:enum('string', 'number', 'bool', 'Json')"`
	Key     string
	Value   string `gorm:"type:text"`
	Comment string `gorm:"type:text"`
}

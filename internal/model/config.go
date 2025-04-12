package model

import (
	"context"

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

	App     App        `gorm:"foreignKey:Name"`
	Type    ConfigType `gorm:"type:enum('string', 'number', 'bool', 'Json')"`
	Key     string
	Value   string `gorm:"type:text"`
	Comment string `gorm:"type:text"`
}

func GetConfig(ctx context.Context, appName, configKey string) (*Config, error) {
	var config Config
	err := db.WithContext(ctx).Where("app = ? AND key = ?", appName, configKey).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func UpdateConfig(ctx context.Context, config *Config) error {
	return db.WithContext(ctx).Save(config).Error
}

package model

import (
	"crypto/rand"
	"encoding/base64"
)

type Settings struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func InitSettings() error {
	if err := db.AutoMigrate(&Settings{}); err != nil {
		return err
	}

	var init_settings []Settings

	// OIDC settings
	init_settings = append(init_settings, Settings{Key: "UseOidc", Value: "false"}) // Whether to use OIDC authentication
	init_settings = append(init_settings, Settings{Key: "OidcProvider", Value: ""})
	init_settings = append(init_settings, Settings{Key: "OidcClientId", Value: ""})
	init_settings = append(init_settings, Settings{Key: "OidcClientSecret", Value: ""})
	init_settings = append(init_settings, Settings{Key: "OidcRedirectUrl", Value: ""})
	init_settings = append(init_settings, Settings{Key: "OidcScopes", Value: "profile"})           // Scopes to request from the OIDC provider
	init_settings = append(init_settings, Settings{Key: "OidcClaimSub", Value: "sub"})             // The claim in the ID token that contains the unique user ID
	init_settings = append(init_settings, Settings{Key: "OidcClaimHasDepartment", Value: "false"}) // Whether the ID token contains a claim that indicates if the user has a department
	init_settings = append(init_settings, Settings{Key: "OidcClaimDepartment", Value: ""})         // The claim in the ID token that contains the user's department

	// JWT settings
	init_settings = append(init_settings, Settings{Key: "JWTValidTime", Value: "24"}) // JWT token valid time in hours
	// Generate a secure random secret key for JWT
	randomSecret := make([]byte, 32)
	if _, err := rand.Read(randomSecret); err == nil {
		init_settings = append(init_settings, Settings{Key: "JWTSecret", Value: base64.StdEncoding.EncodeToString(randomSecret)})
	} else {
		return err
	}

	if err := db.Create(&init_settings).Error; err != nil {
		return err
	}
	return nil
}

func GetSettings(key string) (string, error) {
	var setting Settings
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

func SetSettings(key string, value string) error {
	var setting Settings
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		return err
	}
	setting.Value = value
	return db.Save(&setting).Error
}

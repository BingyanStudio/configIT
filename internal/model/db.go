package model

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func ConnectDB(dsn string, debug bool) (err error) {
	gormLoggerConfig := logger.Config{
		SlowThreshold:        time.Second,
		LogLevel:             logger.Warn,
		ParameterizedQueries: false,
		Colorful:             true,
	}
	if debug {
		gormLoggerConfig.LogLevel = logger.Info
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		gormLoggerConfig,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return errors.Wrap(err, "fail to connect db")
	}

	err = db.AutoMigrate(&App{}, &User{}, &Department{}, &Config{}, &AccessScope{}, &Settings{}, &Namespace{})
	if err != nil {
		return errors.Wrap(err, "fail to auto migrate")
	}
	return nil
}

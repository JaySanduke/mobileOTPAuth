package db

import (
	"mobileOTPAuth/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.UserLoginSession{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

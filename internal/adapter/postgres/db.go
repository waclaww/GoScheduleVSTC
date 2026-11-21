package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetMaxIdleConns(5)

	return db, nil
}

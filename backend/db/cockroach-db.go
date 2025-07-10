package db

import (
    "github.com/SUT-technology/log-analysis/backend/model"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func InitCockroachDB(dsn string) *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to CockroachDB: %v", err)
    }

    if err := db.AutoMigrate(&model.User{}, &model.Project{}); err != nil {
        log.Fatalf("Auto migration failed: %v", err)
    }

    return db
}

package database

import (
	"github.com/Tuxi4k/taskflow/internal/modules/task"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&task.Task{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

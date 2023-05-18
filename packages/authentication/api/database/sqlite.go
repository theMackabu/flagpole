package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLite struct {
	dbName string
}

type IConnection interface {
	GetConnection() (*gorm.DB, error)
}

func (con SQLite) GetConnection() (*gorm.DB, error) {
	con.dbName = os.Getenv("DB_NAME") // move to args
	connection, connectionError := gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if connectionError != nil {
		return nil, connectionError
	}

	return connection, nil
}

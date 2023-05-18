package database

import (
	"flagpole_auth/api/users/models"
	"github.com/gildas/go-logger"
)

type Migrations struct {
	DB IConnection
}

func (migr Migrations) MakeMigrations(log *logger.Logger) {
	connection, conError := migr.DB.GetConnection()
	if conError != nil {
		log.Fatalf("connection error")
	}

	log.Infof("connected to database")
	log.Record("migration", connection.AutoMigrate(&models.User{})).Infof("finished migration")
}

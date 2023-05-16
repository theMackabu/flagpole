package database

import (
	"flagpole_auth/api/users/models"
	"log"
)

type Migrations struct{
	DB IConnection
}

func (migr Migrations) MakeMigrations() {
	connection, conError:= migr.DB.GetConnection()
	if conError != nil{
		log.Println("Connection Error")
	}
	log.Println("Migrating")
	connection.AutoMigrate(&models.User{})
}

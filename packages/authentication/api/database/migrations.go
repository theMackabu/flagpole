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
		log.Println("connection error")
	}
	log.Println("connected to database")
	connection.AutoMigrate(&models.User{})
}

package data

import (
	"log"

	"flagpole_auth/api/database"
	"flagpole_auth/api/users/models"
)

type UserRepository struct {
	DB   database.IConnection
	user *models.User
}

func newRepository() database.IConnection {
	db := UserRepository{
		DB: database.SQLite{},
	}
	return db.DB
}

func (data UserRepository) Find(username string) *models.User {
	connection, connError := newRepository().GetConnection()
	if connError != nil {
		log.Println("connection error")
	}

	connection.Where("username = ?", username).First(&data.user)

	if data.user.Id == "" {
		return data.user
	}

	return data.user
}

func (data UserRepository) Create(user *models.User) (string, error) {
	connection, connError := newRepository().GetConnection()
	if connError != nil {
		log.Println("connection error")
	}
	
	if err := connection.Create(&user).Error; err != nil {		
		return "", err
	} else {
		return user.Id, nil
	}
}

func (data UserRepository) Update(user models.User) {
	connection, connError := data.DB.GetConnection()

	if connError != nil {
		log.Println("connection error")
	}
	connection.Save(&user)
}

func (data UserRepository) Delete(user models.User, id string) {
	connection, connError := data.DB.GetConnection()

	if connError != nil {
		log.Println("connection error")
	}
	connection.Delete(&user, id)
}

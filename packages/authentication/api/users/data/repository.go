package data

import (
	"flagpole_auth/api/database"
	"flagpole_auth/api/users/models"
	"github.com/gildas/go-logger"
)

type UserRepository struct {
	DB   database.IConnection
	user *models.User
	log  *logger.Logger
}

func newRepository() database.IConnection {
	db := UserRepository{
		DB:  database.SQLite{},
		log: logger.Create("authentication", &logger.StdoutStream{Unbuffered: true}),
	}
	return db.DB
}

func (data UserRepository) Find(username string) *models.User {
	connection, connError := newRepository().GetConnection()
	if connError != nil {
		data.log.Errorf("connection error")
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
		data.log.Errorf("connection error")
	}

	if err := connection.Create(&user).Error; err != nil {
		return "", err
	} else {
		return user.Id, nil
	}
}

func (data UserRepository) Update(user models.User) {
	connection, connError := newRepository().GetConnection()

	if connError != nil {
		data.log.Errorf("connection error")
	}
	connection.Save(&user)
}

func (data UserRepository) Delete(user models.User, id string) {
	connection, connError := newRepository().GetConnection()

	if connError != nil {
		data.log.Errorf("connection error")
	}
	connection.Delete(&user, id)
}

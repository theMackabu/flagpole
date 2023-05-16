package data

import "flagpole_auth/api/users/models"

var RepoAccessInterface IUserRepository

type IUserRepository interface {
	Find(id string) *models.User
	Create(user *models.User) (string, error)
}

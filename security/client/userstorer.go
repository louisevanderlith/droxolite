package client

import "github.com/louisevanderlith/droxolite/security/models"

type UserStorer interface {
	FindUser(username string) (string, models.User, error)
	Update(u models.User) error
}

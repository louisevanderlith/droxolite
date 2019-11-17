package impl

import (
	"errors"
	"strconv"

	"github.com/louisevanderlith/droxolite/security/client"
	"github.com/louisevanderlith/droxolite/security/models"
)

type fakeUserStore struct {
	users []models.User
}

func NewFakeUserStore() client.UserStorer {
	result := &fakeUserStore{}

	user, err := models.NewUser("Fake User", "fake@mango.avo", "127.0.0.1", "", "password1")

	if err != nil {
		panic(err)
	}

	user.Verified = true
	result.users = append(result.users, *user)

	return result
}

func (s *fakeUserStore) FindUser(username string) (string, models.User, error) {
	for k, v := range s.users {
		if v.Email == username {
			return strconv.Itoa(k), v, nil
		}
	}

	return "", models.User{}, errors.New("user not found")
}

func (s *fakeUserStore) Update(u models.User) error {
	for _, v := range s.users {
		if v.Email == u.Email {
			v = u
			return nil
		}
	}

	return errors.New("nothing updated")
}

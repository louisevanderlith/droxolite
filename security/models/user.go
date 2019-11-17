package models

import (
	"errors"
	"strings"
	"time"

	"github.com/louisevanderlith/husk"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name        string `hsk:"size(75)"`
	Verified    bool   `hsk:"default(false)"`
	Email       string `hsk:"size(128)"`
	Password    string `hsk:"min(6)"`
	Gravatar    string `hsk:"null"`
	LoginDate   time.Time
	LoginTraces []LoginTrace
	Roles       []Role
	IP          string
	Location    string
}

func (u User) Valid() (bool, error) {
	valid, common := husk.ValidateStruct(&u)

	if !valid {
		return false, common
	}

	if !strings.Contains(u.Email, "@") {
		return false, errors.New("email is invalid")
	}

	return true, nil
}

func NewUser(name, email, ip, location, password string) (*User, error) {
	result := new(User)
	result.Name = name
	result.Email = email
	result.Gravatar = hashGravatar(email)
	result.IP = ip
	result.Location = location
	result.Verified = false

	err := result.SetPassword(password)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *User) SetPassword(password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 11)

	if err != nil {
		return err
	}

	u.Password = string(hashedPwd)

	return nil
}

package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username 	string `json:"username" gorm:"uniqueIndex"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	HashedPassword string `json:"password"`
}

func (user *User) Save() (*User, error) {
	err := Database.Create(user).Error
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func FindUser(id string) (*User, error) {
	var user User
	err := Database.Where("id = ?", id).First(&user).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func FindUserByUsernameOrEmail(usernameOrEmail string) (*User, error) {
	var user User
	err := Database.Where("email = ? OR username = ?", usernameOrEmail, usernameOrEmail).First(&user).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	err := Database.Where("username = ?", username).First(&user).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

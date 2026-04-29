package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username 	string `json:"username" gorm:"uniqueIndex" binding:"required,alphanum,min=3,max=20"`
	Email    string `json:"email" gorm:"uniqueIndex" binding:"required,email"`
	HashedPassword string `json:"password" binding:"required,min=8"`
	Avatar  string `json:"avatar"`
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

func UpdateUserAvatar(id string, avatarUrl string) (*User, error) {
	var user User
	err := Database.Where("id = ?", id).First(&user).Error
	
	if err != nil {
		return &User{}, err
	}

	user.Avatar = avatarUrl
	err = Database.Save(&user).Error

	if err != nil {
		return &User{}, err
	}

	return &user, nil
}
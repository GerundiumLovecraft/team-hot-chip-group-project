package models

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/GerundiumLovecraft/team-hot-chip-group-project/api/src/env"
	"github.com/stretchr/testify/assert"
)

func TestUserSave(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	user := User{
		Username:       "JohnDoe",
		Email:          "johndoe@ntlworld.com",
		HashedPassword: "th15154t3st",
	}
	savedUser, err := user.Save()

	assert.Nil(t, err)
	assert.NotZero(t, savedUser.ID)
	assert.Equal(t, "JohnDoe", savedUser.Username)
	assert.Equal(t, "johndoe@ntlworld.com", savedUser.Email)
	assert.Equal(t, "th15154t3st", savedUser.HashedPassword)
}

func TestFindUser(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	user := User{
		Username:       "JohnDoe",
		Email:          "johndoe@ntlworld.com",
		HashedPassword: "th15154t3st",
	}

	savedUser, err := user.Save()
	assert.Nil(t, err)
	assert.NotZero(t, savedUser.ID)

	idAsString := strconv.Itoa(int(savedUser.ID))
	foundUser, err := FindUser(idAsString)
	assert.Nil(t, err)

	assert.Equal(t, savedUser.ID, foundUser.ID)
	assert.Equal(t, "johndoe@ntlworld.com", foundUser.Email)
	assert.Equal(t, "JohnDoe", foundUser.Username)
}

func TestFindUserByEmail(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	user := User{
		Username:       "JohnDoe",
		Email:          "johndoe@ntlworld.com",
		HashedPassword: "th15154t3st",
	}
	savedUser, err := user.Save()
	assert.Nil(t, err)
	assert.NotZero(t, savedUser.ID)

	foundUser, err := FindUserByUsernameOrEmail(savedUser.Email)
	assert.Nil(t, err)

	assert.Equal(t, "johndoe@ntlworld.com", foundUser.Email)
	assert.Equal(t, "JohnDoe", foundUser.Username)
}

func TestFindUserByUsername(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	user := User{
		Username:       "JohnDoe",
		Email:          "johndoe@ntlworld.com",
		HashedPassword: "th15154t3st",
	}
	savedUser, err := user.Save()
	assert.Nil(t, err)
	assert.NotZero(t, savedUser.ID)
	fmt.Println(savedUser)

	foundUser, err := FindUserByUsernameOrEmail(savedUser.Username)
	assert.Nil(t, err)

	assert.Equal(t, "johndoe@ntlworld.com", foundUser.Email)
	assert.Equal(t, "JohnDoe", foundUser.Username)
}

func TestFindUserByUsernameOrEmailFails(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	_, err := FindUserByUsernameOrEmail("doesnotexist")
	assert.NotNil(t, err)
}

func TestFindByIdFails(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	_, err := FindUser("100001")
	assert.NotNil(t, err)
}

func TestDuplicateEmail(t *testing.T) {
	env.LoadEnv("../../.test.env")
	OpenDatabaseConnection()
	AutoMigrateModels()
	defer Database.Exec("TRUNCATE TABLE users;")

	user1 := User{
		Username:       "JohnDoe",
		Email:          "johndoe@ntlworld.com",
		HashedPassword: "th15154t3st",
	}
	user2 := User{
		Username:       "MichaelDoe",
		Email:          "johndoe@ntlworld.com",
		HashedPassword: "th15164t3st",
	}

	savedUser, err := user1.Save()
	assert.Nil(t, err)

	assert.NotZero(t, savedUser.ID)
	assert.Equal(t, "JohnDoe", savedUser.Username)
	assert.Equal(t, "johndoe@ntlworld.com", savedUser.Email)
	assert.Equal(t, "th15154t3st", savedUser.HashedPassword)

	_, err = user2.Save()
	assert.NotNil(t, err)
}

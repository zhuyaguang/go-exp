package model

import (
	"errors"
	"fmt"
	"go-zero-api/service/internal/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)



// CreateUser add a new user
func CreateUser(name, passWord string) interface{} {
	//open a db connection
	db.Create(&types.RegisterRequest{
		Username:    "name",
		Password:    "passWord",
		Phonenumber: "123",
	})
	return nil
}

// UserLogin login system
func UserLogin(name, passWord string) error {
	var user types.LoginRequest
	db.Select("password").Where("username = ?", name).First(&user)
	if user.Password != passWord {
		return errors.New("password is not correct")
	}
	//todo println("add 认证 鉴权 逻辑")
	return nil
}

// GetUser get a single user
func GetUser(userID string) (result types.RegisterRequest, err error) {
	var user types.RegisterRequest
	db.First(&user, userID)

	resultUser := types.RegisterRequest{Username: user.Username, Password: user.Password}
	return resultUser, nil
}

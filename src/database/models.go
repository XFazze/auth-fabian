package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       uint
	Username string
	Email    string
	Password string
}

type User_tokens struct {
	gorm.Model
	Token string
	Uid   uint
}

type Forgot_password_code struct {
	gorm.Model
	Code  string
	Email string
}

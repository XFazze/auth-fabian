package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type User_tokens struct {
	gorm.Model
	Token    string
	Username string
}

type Forgot_password_code struct {
	gorm.Model
	Code  string
	Email string
}

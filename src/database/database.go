package database

import (
	"errors"
	"fmt"
	"os"
	"auth-fabian/src/base"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup_migrate_db() {
	_, err := os.Stat(os.Getenv("DATABASE_FILE"))
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("database does not exist")
		fmt.Println(os.Getenv("DATABASE_FILE"))
		file, err := os.Create(os.Getenv("DATABASE_FILE")) // FIXME does not create instant?
		base.CheckErr(err)
		file.Close()
	}

	db := OpenDB()
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &User_tokens{}, &Forgot_password_code{})
}

func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_FILE")), &gorm.Config{})
	base.CheckErr(err)
	return db
}

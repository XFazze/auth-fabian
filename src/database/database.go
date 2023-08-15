package database

import (
	"auth-fabian/src/base"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup_migrate_db() {
	_, err := os.Stat(os.Getenv("DATABASE_FILE"))
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("database does not exist")
		file, err := os.Create(os.Getenv("DATABASE_FILE")) // FIXME does not create instant?
		base.CheckErr(err)
		file.Close()
	}

	db := OpenDB()
	db.AutoMigrate(&User{}, &User_tokens{}, &Forgot_password_code{})
}

func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_FILE")), &gorm.Config{})
	base.CheckErr(err)
	return db
}

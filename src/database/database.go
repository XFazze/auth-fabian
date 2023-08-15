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
		base.Check_err(err)
		file.Close()
	}

	db := Open_DB()
	db.AutoMigrate(&User{}, &User_tokens{}, &Forgot_password_code{})
}

func Open_DB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_FILE")), &gorm.Config{})
	base.Check_err(err)
	return db
}

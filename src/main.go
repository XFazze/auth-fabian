package main

import (
	"os"

	"auth-fabian/src/database"
	"auth-fabian/src/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()               // default .env file
	godotenv.Load(".env_secrets") // .env_secrets file

	gin.SetMode(os.Getenv("GIN_MODE"))
	database.Setup_migrate_db()

}
func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"localhost"})
	r.LoadHTMLGlob("templates/**/*") // TODO load templates folder not only subdirectories

	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	r.GET("/", routes.Index)

	r.GET("/login", routes.Login)
	r.POST("/login", routes.Login_form)

	r.GET("/logout", routes.Logout)

	r.GET("/signup", routes.Signup)
	r.POST("/signup", routes.Signup_form)

	r.GET("/forgot_password", routes.Forgot_password)
	r.POST("/forgot_password", routes.Forgot_password_form)
	r.GET("/forgot_password/:email/:code", routes.Forgot_password_change)
	r.POST("/forgot_password/:email/:code", routes.Forgot_password_code_form)

	r.GET("/validate_token/:token", routes.Validate_token)

	r.Run()

}

/*
func login(username_mail, passwd string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(username_mail, "@") {
		db.Exec(create)
	}
	if err != nil {
		fmt.Println(err)
	}

}
*/

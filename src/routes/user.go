package routes

import (
	"auth-fabian/src/base"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"auth-fabian/src/database"
	"auth-fabian/src/external_api"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
)

var (
	verifier = emailverifier.NewVerifier().EnableSMTPCheck()
)

func get_cookies(c *gin.Context) (string, uint) {
	cookie, err := c.Cookie("user_token")
	if err != nil {
		cookie = ""
	}

	uint_id := uint64(0)
	id, err := c.Cookie("id")
	if err == nil {
		uint_id, err = strconv.ParseUint(id, 10, 32)
		base.Check_err(err)
	}

	return cookie, uint(uint_id)
}
func Index(c *gin.Context) {
	cookie, id := get_cookies(c)
	logged_in, _, _ := check_login(c, cookie, id)
	var user database.User
	if logged_in {
		db := database.Open_DB()
		db.Find(&user, "id = ?", id)
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"logged_in": logged_in, "username": user.Username})
}

func Redirect_user(token string, expire_seconds string, id uint, c *gin.Context) {
	redirect_url := c.DefaultQuery("redirect", "/")
	if redirect_url == "/" {
		c.Redirect(http.StatusFound, "/")
		return
	}
	fmt.Println("redirect url", redirect_url)
	params := url.Values{}
	params.Set("token", token)
	params.Set("expire_seconds", expire_seconds)
	params.Set("id", strconv.FormatUint(uint64(id), 10))
	location, err := url.Parse(redirect_url)
	base.Check_err(err)
	if location.Host == "" {
		location, err = url.Parse("https://" + redirect_url)
		base.Check_err(err)
	}
	location.RawQuery = params.Encode()

	fmt.Println("host", location.Host)
	fmt.Println("path", location.Path)
	fmt.Println("RawQuery", location.RawQuery)
	fmt.Println("RequestURI", location.RequestURI())

	c.Redirect(http.StatusFound, location.String())
}

func Login(c *gin.Context) {
	cookie, id := get_cookies(c)
	logged_in, token, expire_seconds := check_login(c, cookie, id)
	if logged_in {
		Redirect_user(token, expire_seconds, id, c)

	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	}
}

func Login_form(c *gin.Context) {
	db := database.Open_DB()
	user_identifier := c.PostForm("user_identifier")
	password := c.PostForm("password")
	var user database.User
	var id_type string
	if strings.Contains(user_identifier, "@") {
		db.Find(&user, "email = ?", user_identifier)
		id_type = "email"
	} else {
		db.Find(&user, "username = ?", user_identifier)
		id_type = "username"
	}

	if user.Username == "" { // user doesnt exist
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": fmt.Sprintf("No user with %s: %s", id_type, user_identifier),
		})
		return

	}
	if base.Check_password_hash(password, user.Password) {
		login_user(c, user.Id)

	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "Wrong credentials.",
		})
	}

}

func Logout(c *gin.Context) {
	c.SetCookie("id", "", -1, "/", os.Getenv("DOMAIN"), true, true)
	c.SetCookie("username", "", -1, "/", os.Getenv("DOMAIN"), true, true)
	c.Redirect(http.StatusFound, "/")
}
func Signup(c *gin.Context) {
	cookie, id := get_cookies(c)
	logged_in, token, expire_seconds := check_login(c, cookie, id)
	if logged_in {
		Redirect_user(token, expire_seconds, id, c)

	} else {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{})
	}

}

func Signup_form(c *gin.Context) {
	db := database.Open_DB()

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	ret, err := verifier.Verify(email)
	if err != nil {
		fmt.Println("verify email address failed, error is: ", err)
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"error":          "Not a valid email.(moron)",
			"reset_password": false,
		})
		return
	}
	if !ret.Syntax.Valid {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"error":          "Not a valid email.",
			"reset_password": false,
		})
		return
	}

	var user database.User
	result := db.Find(&user, "username = ?", username)
	if result.RowsAffected != 0 { // if username exists
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"error":          fmt.Sprintf(`Username: "%s" is already taken.`, username),
			"reset_password": false,
		})
		return

	}
	result = db.Find(&user, "email = ?", email)
	if result.RowsAffected != 0 { // if email exists
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"error":          fmt.Sprintf(`Email: "%s" is already registered.`, email),
			"reset_password": true,
		})
		return

	}
	hased_password := base.Hash_password(password)

	db.Create(&database.User{Username: username, Email: email, Password: hased_password})
	db.Find(&user, "username = ?", username)

	login_user(c, user.Id)
}

func login_user(c *gin.Context, id uint) {
	db := database.Open_DB()
	generated_token := base.Generate_secure_token(32)
	db.Create(&database.User_tokens{Token: generated_token, Uid: id})
	expiry_seconds, err := strconv.Atoi(os.Getenv("EXPIRY_SECONDS"))
	base.Check_err(err)
	c.SetCookie("user_token", generated_token, expiry_seconds, "/", os.Getenv("DOMAIN"), true, true)
	c.SetCookie("id", strconv.FormatUint(uint64(id), 10), expiry_seconds, "/", os.Getenv("DOMAIN"), true, true)
	Redirect_user(generated_token, os.Getenv("EXPIRY_SECONDS"), id, c)

}

func check_login(c *gin.Context, token string, id uint) (bool, string, string) {

	db := database.Open_DB()
	var user_token database.User_tokens
	db.Find(&user_token, "token = ? AND uid = ?", token, id)

	expiry_seconds, err := strconv.Atoi(os.Getenv("EXPIRY_SECONDS"))
	base.Check_err(err)
	if user_token.Token != "" && user_token.UpdatedAt.Before(time.Now().Add(time.Duration(expiry_seconds*10e8))) {
		time_difference := user_token.UpdatedAt.Sub(time.Now())
		return true, user_token.Token, fmt.Sprintf("%f", time_difference.Seconds())
	}
	return false, "", ""

}

func Forgot_password(c *gin.Context) {
	cookie, id := get_cookies(c)
	logged_in, token, expire_seconds := check_login(c, cookie, id)
	if logged_in {
		Redirect_user(token, expire_seconds, id, c)

	} else {
		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{})
	}

}

func Forgot_password_form(c *gin.Context) {
	db := database.Open_DB()
	email_input := c.PostForm("email")

	ret, err := verifier.Verify(email_input)
	if err != nil {
		fmt.Println("verify email address failed, error is: ", err)
		return
	}
	if !ret.Syntax.Valid {
		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
			"error":          "Not a valid email.",
			"reset_password": false,
			"done":           false,
		})
		return
	}

	var user database.User
	result := db.Find(&user, "email = ?", email_input)

	if os.Getenv("EMAIL_LOGIN") == "" || os.Getenv("EMAIL_PASSWORD") == "" || os.Getenv("EMAIL_PUBLIC") == "" {
		log.Fatal("EMAIL_LOGIN, EMAIL_PUBLIC or EMAIL_PASSWORD is not set")
	}
	if result.RowsAffected != 0 { // if email exists
		code := base.Generate_secure_token(30)
		db.Create(&database.Forgot_password_code{Code: code, Email: email_input})
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		link := fmt.Sprintf("%s://%s/forgot_password/%s/%s", scheme, c.Request.Host, email_input, code)
		external_api.Send_mail(email_input, "Password reset", fmt.Sprintf(`<h1>Password request</h1>
		<p>A password reset was requested on this email. 
		You can ignore this if it wasnt requested by you.  </p><br />
		<a href="% s">Click to reset password.</a>
		<p>Or paste the link into youre browser. Link:%s</p>`, link, link))

	}
	c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
		"sent":  true,
		"email": email_input,
	})

}

func Forgot_password_change(c *gin.Context) {
	db := database.Open_DB()
	code := c.Param("code")
	email := c.Param("email")
	var code_datbase database.Forgot_password_code
	db.Find(&code_datbase, "code = ? AND email =?", code, email)
	if code_datbase.Code == "" {
		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
			"code_tab":   true,
			"code_error": "Password code is not valid. Might be expired.",
		})
	} else {

		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
			"code_tab": true,
		})
	}

}
func Forgot_password_code_form(c *gin.Context) {
	db := database.Open_DB()
	code := c.Param("code")
	email := c.Param("email")
	var code_datbase database.Forgot_password_code
	db.Find(&code_datbase, "code = ? AND email =?", code, email)
	if code_datbase.Code == "" {
		c.HTML(http.StatusOK, "forgot_password.tmpl", gin.H{
			"code_tab":   true,
			"code_error": "Password code is not valid. Might be expired.",
		})
	} else {
		password := c.PostForm("password")
		db.Model(&database.User{}).Where("email = ?", email).Update("password", base.Hash_password(password))

		var user database.User
		db.Find(&user, "email = ?", email)
		login_user(c, user.Id)
	}

}

func Validate_token(c *gin.Context) {
	token := c.Param("token")
	id := c.Param("id")
	id_str, err := strconv.ParseUint(id, 10, 32)
	base.Check_err(err)
	logged_in, _, _ := check_login(c, token, uint(id_str))
	db := database.Open_DB()
	var user database.User
	db.Find(&user, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"valid":    logged_in,
		"username": user.Username,
	})

}
func Delete_account(c *gin.Context) {
	c.HTML(http.StatusOK, "delete_account.tmpl", gin.H{})

}

func Delete_account_form(c *gin.Context) {
	db := database.Open_DB()

	user_identifier := c.PostForm("user_identifier")
	password := c.PostForm("password")
	var user database.User
	var id_type string
	if strings.Contains(user_identifier, "@") {
		db.Find(&user, "email = ?", user_identifier)
		id_type = "email"
	} else {
		db.Find(&user, "username = ?", user_identifier)
		id_type = "username"
	}

	if user.Username == "" { // user doesnt exist
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": fmt.Sprintf("No user with %s: %s", id_type, user_identifier),
		})
		return

	}
	if base.Check_password_hash(password, user.Password) {
		db.Delete(&user)
		var user_token database.User_tokens
		db.Delete(&user_token, "uid = ?", user.Id)

	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error": "Wrong credentials.",
		})
	}
	c.Redirect(http.StatusFound, "/")

}

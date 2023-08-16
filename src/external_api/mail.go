package external_api

import (
	"auth-fabian/src/base"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func Send_mail(to, subject, html string) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Auth fabian Support <%s>", os.Getenv("EMAIL_PUBLIC"))
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", os.Getenv("EMAIL_LOGIN"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com"))
	base.Check_err(err)

}

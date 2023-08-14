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
	e.From = fmt.Sprintf("Auth fabian suppoer <%s>", os.Getenv("EMAIL"))
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)
	fmt.Println(os.Getenv("EMAIL"), os.Getenv("EMAIL_PASSWORD"))
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", os.Getenv("EMAIL"), os.Getenv("EMAIL_PASSWORD"), "smtp.gmail.com"))
	base.CheckErr(err)

}

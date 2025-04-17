package config

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

var mail *gomail.Dialer

func linkService() {
	mail = gomail.NewDialer(
		viper.GetString("EMAIL_HOST"),
		viper.GetInt("EMAIL_PORT"),
		viper.GetString("EMAIL_USERNAME"),
		viper.GetString("EMAIL_PASSWORD"))
}

func SendEmail(email, formName, subject, text string) error {
	if mail == nil {
		linkService()
	}

	tmail := viper.GetString("EMAIL_USERNAME")

	m := gomail.NewMessage()
	m.SetAddressHeader("From", tmail, formName)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)
	return mail.DialAndSend(m)
}

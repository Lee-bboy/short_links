package common

import (
	"shortlinks/config"
	"shortlinks/mailer"
)

var Mailer *mailer.Mailer

func init() {
	address := config.GetConf("mailer.host", "")
	tls := config.GetConf("mailer.ssl", "0") == "1"
	username := config.GetConf("mailer.username", "")
	password := config.GetConf("mailer.password", "")

	var err error
	Mailer, err = mailer.NewMailer(address, username, password, tls)
	if err != nil {
		panic(err)
	}
}

func MailTo(to string) *mailer.Mail {
	return mailer.NewMail(config.GetConf("mailer.username", ""), to)
}

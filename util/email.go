package util

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"os"
	"text/template"
)

func SendEmail(receiver, username, password string) error {
	sender := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PASS")

	t, err := template.ParseFiles("util/template.html")

	if err != nil {
		return err
	}

	var body bytes.Buffer

	err = t.Execute(&body, struct {
		Username string
		Password string
	}{
		Username: username,
		Password: password,
	})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "سامانه تردد مهمان")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, sender, pass)

	//d.TLSConfig = &tls.Config{InsecureSkipVerify: false}

	return d.DialAndSend(m)
}

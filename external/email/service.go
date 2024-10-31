package sendemail

import (
	"bytes"
	"html/template"
	"net/smtp"
)

var auth smtp.Auth

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

type AuthEmail struct {
	Email    string
	Password string
	Smtp     string
}

func NewRequest(to []string, subject, body string, email AuthEmail) *Request {

	auth = smtp.PlainAuth("", email.Email, email.Password, email.Smtp)
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "eko@thegreatsoft.com", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

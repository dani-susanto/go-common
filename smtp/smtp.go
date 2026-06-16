package smtp

import (
	gomail "gopkg.in/gomail.v2"
)

type SMTP interface {
	Send(subject string, to string, content string) error
}

type smtp struct {
	host     string
	port     int
	username string
	password string
	from     string
	tls      bool
}

func New(
	host string,
	port int,
	username string,
	password string,
	from string,
	tls bool,
) SMTP {
	return &smtp{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		tls:      tls,
	}
}

func (s *smtp) Send(subject string, to string, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	dialer := gomail.NewDialer(
		s.host,
		s.port,
		s.username,
		s.password,
	)
	// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	dialer.SSL = s.tls

	return dialer.DialAndSend(msg)
}

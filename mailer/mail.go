package mailer

import (
	"fmt"
	"time"
)

type Mail struct {
	errorCh chan error

	from    string
	to      []string
	subject []byte
	body    []byte
}

func NewMail(from, to string) *Mail {
	mail := &Mail{
		errorCh: make(chan error),
		from:    from,
		to:      make([]string, 0),
	}
	mail.to = append(mail.to, to)

	return mail
}

func (m *Mail) Subject(subject []byte) *Mail {
	m.subject = subject
	return m
}

func (m *Mail) SubjectFromString(subject string) *Mail {
	m.subject = []byte(subject)
	return m
}

func (m *Mail) Body(body []byte) *Mail {
	m.body = body
	return m
}

func (m *Mail) BodyFromString(body string) *Mail {
	m.body = []byte(body)
	return m
}

func (m *Mail) String() string {
	var data string

	//mail header
	data += fmt.Sprintf("FROM:%s\r\n", m.from)
	data += fmt.Sprintf("TO:%s\r\n", m.to[0])
	data += fmt.Sprintf("SUBJECT:%s\r\n", m.subject)
	data += fmt.Sprintf("Date:%s\r\n", time.Now().Format(time.RFC1123))
	data += "MIME-Version:1.0\r\n"
	data += "Content-Type:text/plain\r\n\r\n"

	//mail body
	data += string(m.body) + "\r\n"

	return data
}

package mailer

import (
	ttls "crypto/tls"
	"io"
	"net"
	"net/smtp"
	"sync"
)

type plainAuthOverTLSConn struct {
	smtp.Auth
}

func PlainAuthOverTLSConn(identity, username, password, host string) smtp.Auth {
	return &plainAuthOverTLSConn{smtp.PlainAuth(identity, username, password, host)}
}

func (a *plainAuthOverTLSConn) Start(server *smtp.ServerInfo) (string, []byte, error) {
	server.TLS = true
	return a.Auth.Start(server)
}

type Mailer struct {
	mailCh chan *Mail

	host     string
	port     string
	username string
	password string
	tls      bool //是否使用了TLS加密

	client *smtp.Client

	running      bool
	runningMutex *sync.Mutex
}

func NewMailer(address, username, password string, tls bool) (*Mailer, error) {
	var client *smtp.Client
	var err error

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", username, password, host)

	if tls {
		var tlsConn *ttls.Conn
		tlsConn, err = ttls.Dial("tcp", address, nil)
		if err != nil {
			return nil, err
		}

		client, err = smtp.NewClient(tlsConn, host)
		if err != nil {
			return nil, err
		}
	} else {
		var conn net.Conn
		conn, err = net.Dial("tcp", address)
		if err != nil {
			return nil, err
		}

		client, err = smtp.NewClient(conn, host)
		if err != nil {
			return nil, err
		}
	}

	//auth
	err = client.Auth(&plainAuthOverTLSConn{auth})
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		mailCh:       make(chan *Mail),
		host:         host,
		port:         port,
		username:     username,
		password:     password,
		tls:          tls,
		client:       client,
		running:      true,
		runningMutex: new(sync.Mutex),
	}

	//运行
	go mailer.run()

	return mailer, nil
}

func (m *Mailer) Send(mail *Mail) error {
	m.mailCh <- mail

	return <-mail.errorCh
}

func (m *Mailer) Close() {
	m.runningMutex.Lock()
	defer m.runningMutex.Unlock()

	if !m.running {
		return
	}

	m.running = false
	close(m.mailCh)
}

func (m *Mailer) run() {
	var mail *Mail
	var err error

	for mail = range m.mailCh {
		err = m.client.Mail(mail.from)
		if err != nil {
			mail.errorCh <- err
			continue
		}

		err = m.client.Rcpt(mail.to[0])
		if err != nil {
			mail.errorCh <- err
			continue
		}

		wr, err := m.client.Data()
		if err != nil {
			mail.errorCh <- err
			continue
		}
		_, err = io.WriteString(wr, mail.String())
		if err != nil {
			mail.errorCh <- err
			continue
		}
		wr.Close()

		mail.errorCh <- nil
	}

	m.client.Quit()
}

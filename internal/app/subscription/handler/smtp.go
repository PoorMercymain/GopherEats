package handler

import (
	"context"
	"encoding/base64"
	"net/smtp"
	"strings"
)

type smtpSender struct {
	username string
	password string
	server   string
	port     string
}

func (s *smtpSender) SendEmail(ctx context.Context, to string, subject string, message string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.server)

	var strBuilder strings.Builder

	msgStrSlice := []string{"From: ", s.username, "\r\nTo: ", to, "\r\nSubject: ", subject,
		"\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n", "Content-Transfer-Encoding: base64\r\n", "\r\n\r\n", base64.StdEncoding.EncodeToString([]byte(message))}

	for _, str := range msgStrSlice {
		_, err := strBuilder.WriteString(str)
		if err != nil {
			return err
		}
	}

	err := smtp.SendMail(s.server+":"+s.port, auth, s.username, []string{to}, []byte(strBuilder.String()))
	if err != nil {
		return err
	}

	return nil
}

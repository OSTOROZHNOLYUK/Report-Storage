package notifications

import (
	"fmt"
	"net/smtp"
)

const (
	// Тема письма об изменении статуса заявки.
	statusSubject = "Статус заявки изменен"
	// Тело письма об изменении статуса заявки.
	statusBody = "Статус Вашей заявки в проекте \"Осторожно, люк!\" изменен.\nНовый статус: "
	// Тема письма о создании новой заявки.
	newSubject = "Создана новая заявка"
	// Тело письма о создании новой заявки.
	newBody = "Создана новая заявка в проекте \"Осторожно, люк!\".\nЗаявка отобразится на карте после проверки модератором."
)

// SMTP - структура клиента SMTP сервера.
type SMTP struct {
	sender   string
	login    string
	password string
	host     string
	port     string
}

// New - конструктор клиента SMTP сервера.
func New(sender, login, passwd, host, port string) *SMTP {
	return &SMTP{
		sender:   sender,
		login:    login,
		password: passwd,
		host:     host,
		port:     port,
	}
}

// StatusChanged отправляет уведомление на почту target об изменении
// статуса на status.
func StatusChanged(mail *SMTP, target, status string) error {
	auth := smtp.PlainAuth("", mail.login, mail.password, mail.host)

	body := statusBody + status
	msg := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s\r\n", target, statusSubject, body,
	)
	addr := fmt.Sprintf("%s:%s", mail.host, mail.port)

	err := smtp.SendMail(addr, auth, mail.sender, []string{target}, []byte(msg))
	return err
}

// NewReport отправляет уведомление на почту target о создании новой заявки.
func NewReport(mail *SMTP, target string) error {
	auth := smtp.PlainAuth("", mail.login, mail.password, mail.host)

	msg := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s\r\n", target, newSubject, newBody,
	)
	addr := fmt.Sprintf("%s:%s", mail.host, mail.port)

	err := smtp.SendMail(addr, auth, mail.sender, []string{target}, []byte(msg))
	return err
}

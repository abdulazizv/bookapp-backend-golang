package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"

	"gitlab.com/bookapp/config"
)

func MailSender(mail string, code string) error {
	cfg := config.Load()
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", cfg.MailUsername, cfg.MailPassword, cfg.SmtpHost)

	var emailBody bytes.Buffer
	data := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}
	tmp, err := template.ParseFiles("./api/handler/v1/html/index.html")
	if err != nil {
		return err
	}

	if err := tmp.Execute(&emailBody, &data); err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}

	// Message data
	from := "noreply@gmail.com"
	subject := "Subject: Chat App\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	to := []string{mail}

	message := []byte(subject + mime + emailBody.String())

	// Connect to the server and send message
	smtpUrl := cfg.SmtpHost + ":587"
	err = smtp.SendMail(smtpUrl, auth, from, to, message)
	if err != nil {
		return err
	}
	return err
}

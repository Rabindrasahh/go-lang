package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func parseTemplate(templateName string, data interface{}) (string, error) {
	templatePath := filepath.Join("view", "email", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func initSMTPConfig() (string, string, string, string) {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	smtpUser := getEnv("SMTP_USER", "default_user")
	smtpPassword := getEnv("SMTP_PASSWORD", "default_password")
	smtpHost := getEnv("SMTP_HOST", "localhost")
	smtpPort := getEnv("SMTP_PORT", "25")

	return smtpUser, smtpPassword, smtpHost, smtpPort
}

func SendHTMLEmail(to, subject, templateName string, data interface{}) error {
	smtpUser, smtpPassword, smtpHost, smtpPort := initSMTPConfig()
	body, err := parseTemplate(templateName, data)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s",
		smtpUser, to, subject, body)

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(msg))
}

func SendVerificationEmail(to, subject, verificationURL string) error {
	data := map[string]interface{}{
		"VerificationURL": verificationURL,
	}
	return SendHTMLEmail(to, subject, "verification_email.html", data)
}

func SendPasswordChangeNotification(to, name string) error {
	data := map[string]interface{}{
		"Name": name,
	}
	return SendHTMLEmail(to, "Password Change Notification", "password_change_email.html", data)
}

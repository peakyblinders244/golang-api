package helper

import (
	"bytes"
	"crypto/tls"
	"github.com/joho/godotenv"
	"github.com/k3a/html2text"
	"golang-api/config"
	"golang-api/entity"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type EmailData struct {
	URL       string
	FirstName string
	Code      string
	Subject   string
}

// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *entity.User, data *EmailData) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Sender data.
	env := config.LoadEnv()
	from := env.EMAIL_FROM
	smtpPass := env.SMTP_PASS
	smtpUser := env.SMTP_USER
	to := user.Email
	smtpHost := env.SMTP_HOST
	smtpPort, err := strconv.Atoi(env.SMTP_PORT)
	if err != nil {
		log.Fatal("failed to parse port integer", err)
	}

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}

}

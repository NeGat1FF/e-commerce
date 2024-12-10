package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	// Mailtrap
	MailtrapURL = "https://send.api.mailtrap.io/api/send"
)

type EmailType string

const (
	EmailConfirmation EmailType = "email_confirm"
	PasswordReset     EmailType = "password_reset"
)

type EmailSender struct {
	jwt    string
	client *http.Client

	templates map[string]*template.Template
}

type FromPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ToPayload struct {
	Email string `json:"email"`
}

type EmailPayload struct {
	From     FromPayload `json:"from"`
	To       []ToPayload `json:"to"`
	Subject  string      `json:"subject"`
	Text     string      `json:"text"`
	HTML     string      `json:"html"`
	Category string      `json:"category"`
}

func NewEmailPayload(toEmail []string, subject, text, html, category string) (*bytes.Reader, error) {
	// Convert toEmail to ToPayload
	to := make([]ToPayload, 0)
	for _, email := range toEmail {
		to = append(to, ToPayload{Email: email})
	}

	p := EmailPayload{
		From: FromPayload{
			Email: "hello@demomailtrap.com",
			Name:  "Mailtrap Test",
		},
		To:       to,
		Subject:  subject,
		Text:     text,
		HTML:     html,
		Category: category,
	}

	data, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func NewEmailSender(client *http.Client, jwt string) *EmailSender {
	return &EmailSender{
		client: client,
		jwt:    jwt,
	}
}

func (m *EmailSender) InitTemplates(path string) error {
	// Load all templates from the path
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	m.templates = make(map[string]*template.Template)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".html") {
			tmpl, err := template.ParseFiles(path + file.Name())
			if err != nil {
				return err
			}
			m.templates[strings.TrimSuffix(file.Name(), ".html")] = tmpl
		}
	}

	return nil
}

func (m *EmailSender) SendMail(to []string, tp EmailType, data map[string]string) error {
	// Get the template
	tmpl, ok := m.templates[string(tp)]
	if !ok {
		return fmt.Errorf("template %s not found", tp)
	}

	// Execute the template
	var text bytes.Buffer
	if err := tmpl.Execute(&text, data); err != nil {
		return err
	}

	var subject string

	switch tp {
	case EmailConfirmation:
		subject = "Email Confirmation"
	case PasswordReset:
		subject = "Password Reset"
	}

	payload, err := NewEmailPayload(to, subject, text.String(), text.String(), "category")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", MailtrapURL, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+m.jwt)
	req.Header.Add("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

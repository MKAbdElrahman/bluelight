package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bluelight.mkcodedev.com/src/infrastructure/mailer/templates"
)

type Mailer struct {
	client *http.Client
	config Config
}

type Config struct {
	URL    string
	APIKey string
	Sender EmailAddress
}

type EmailPayload struct {
	From     EmailAddress   `json:"from"`
	To       []EmailAddress `json:"to"`
	Subject  string         `json:"subject"`
	HTML     string         `json:"html"`
	Category string         `json:"category"`
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// NewMailer initializes a new Mailer with the given configuration.
func NewMailer(config Config) *Mailer {
	return &Mailer{
		client: &http.Client{},
		config: config,
	}
}

// WelcomeNewRegisteredUser sends a welcome email to a new user.
func (m *Mailer) WelcomeNewRegisteredUser(ctx context.Context, recipientEmail, recipientName string) error {
	body, err := m.renderTemplate(ctx, recipientName)
	if err != nil {
		return err
	}
	payload := m.buildEmailPayload(recipientEmail, recipientName, body)
	return m.sendEmail(payload)
}


// renderTemplate renders the welcome email template.
func (m *Mailer) renderTemplate(ctx context.Context, recipientName string) (string, error) {
	var body bytes.Buffer
	err := templates.UserWelcome(recipientName).Render(ctx, &body)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}
	return body.String(), nil
}


// buildEmailPayload constructs the email payload for the welcome email.
func (m *Mailer) buildEmailPayload(recipientEmail, recipientName, body string) EmailPayload {
	return EmailPayload{
		From: m.config.Sender,
		To: []EmailAddress{
			{
				Email: recipientEmail,
				Name:  recipientName,
			},
		},
		Subject:  "Welcome to Bluelight, " + recipientName + "!",
		HTML:     body,
		Category: "User Registration",
	}
}

// sendEmail sends the constructed email payload using the configured email service.
func (m *Mailer) sendEmail(payload EmailPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", m.config.URL, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.config.APIKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to send email, status: %d, response: %s", res.StatusCode, string(body))
	}

	return nil
}

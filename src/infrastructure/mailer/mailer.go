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
	url    string
	apiKey string
	sender string
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

func New(url, apiKey, sender string) Mailer {
	client := &http.Client{}
	return Mailer{
		client: client,
		url:    url,
		apiKey: apiKey,
		sender: sender,
	}
}

func (m Mailer) WelcomeNewRegisteredUser(recipientEmail, recipientName string) error {
	method := "POST"

	var body bytes.Buffer
	err := templates.UserWelcome(recipientName).Render(context.Background(), &body)

	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	payload := EmailPayload{
		From: EmailAddress{
			Email: m.sender,
			Name:  "Bluelight",
		},
		To: []EmailAddress{
			{
				Email: recipientEmail,
				Name:  recipientName,
			},
		},
		Subject:  "Welcome to Bluelight, " + recipientName + "!",
		HTML:     body.String(),
		Category: "User Registration",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest(method, m.url, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))
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

package mailjet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/manciniraka/urbioxe/internal/constant"
	"github.com/manciniraka/urbioxe/internal/logger"
)

type Config struct {
	BaseURL     string
	APIKey      string
	SecretKey   string
	SenderEmail string
	SenderName  string
}

type Client struct {
	config     Config
	httpClient *http.Client
}

func New(cfg Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type payloadSendEmail struct {
	Messages []message `json:"Messages"`
}

type message struct {
	From     from   `json:"From"`
	To       []to   `json:"To"`
	Subject  string `json:"Subject"`
	TextPart string `json:"TextPart"`
	HTMLPart string `json:"HTMLPart"`
}

type from struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type to struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

func (c *Client) endpoint() string {
	return c.config.BaseURL + "/v3.1/send"
}

func (c *Client) SendEmail(
	toEmail string,
	toName string,
	subject string,
	textBody string,
	htmlBody string,
) error {

	var lastErr error

	for attempt := 1; attempt <= constant.DefaultMaxRetry; attempt++ {

		lastErr = c.send(
			toEmail,
			toName,
			subject,
			textBody,
			htmlBody,
		)

		if lastErr == nil {

			if attempt > 1 {

				logger.Log.Info(
					"mailjet email sent after retry",
					"tag", constant.LogTagMailjet,
					"attempt", attempt,
					"recipient", toEmail,
				)

			}

			return nil
		}

		logger.Log.Warn(
			"mailjet send failed",
			"tag", constant.LogTagMailjet,
			"attempt", attempt,
			"recipient", toEmail,
			"error", lastErr,
		)

		if attempt < constant.DefaultMaxRetry {
			time.Sleep(
				constant.DefaultRetryDelay,
			)
		}
	}

	logger.Log.Error(
		"mailjet failed after retry",
		"tag", constant.LogTagMailjet,
		"recipient", toEmail,
		"error", lastErr,
	)

	return lastErr
}

func (c *Client) send(
	toEmail string,
	toName string,
	subject string,
	textBody string,
	htmlBody string,
) error {

	payload := payloadSendEmail{
		Messages: []message{
			{
				From: from{
					Email: c.config.SenderEmail,
					Name:  c.config.SenderName,
				},
				To: []to{
					{
						Email: toEmail,
						Name:  toName,
					},
				},
				Subject:  subject,
				TextPart: textBody,
				HTMLPart: htmlBody,
			},
		},
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.endpoint(),
		bytes.NewBuffer(payloadByte),
	)
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(
			c.config.APIKey +
				":" +
				c.config.SecretKey,
		),
	)

	req.Header.Set(
		"Authorization",
		"Basic "+auth,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf(
			"mailjet returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}

func (c *Client) SendWelcomeEmail(
	name string,
	email string,
) error {

	subject := "Welcome to Urbioxe!"

	textBody := fmt.Sprintf(
		`Welcome, %s!

Thank you for registering at Urbioxe.

Your account has been created successfully.

You can now:
- Report city issues
- Track your reports
- Stay updated with district and city information

Regards,
Urbioxe Team`,
		name,
	)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
</head>

<body style="font-family: Arial, Helvetica, sans-serif; line-height:1.6; color:#333333;">

<h2>👋 Welcome, %s!</h2>

<p>
Thank you for registering at <strong>Urbioxe</strong>.
</p>

<p>
Your account has been created successfully.
</p>

<p>
You can now:
</p>

<ul>
<li>📍 Report city issues</li>
<li>📋 Track your reports</li>
<li>🌤️ Stay updated with district and city information</li>
</ul>

<p>
We're excited to have you as part of our community.
</p>

<br>

<p>
Regards,<br>
<strong>Urbioxe Team</strong>
</p>

</body>
</html>`,
		name,
	)

	return c.SendEmail(
		email,
		name,
		subject,
		textBody,
		htmlBody,
	)
}

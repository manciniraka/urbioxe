package mailjet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	Messages []message `json:"messages"`
}

type message struct {
	From     from   `json:"from"`
	To       []to   `json:"to"`
	Subject  string `json:"subject"`
	TextPart string `json:"text_part"`
	HTMLPart string `json:"HTML_part"`
}

type from struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type to struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (c *Client) SendEmail(
	toName string,
	toEmail string,
	subject string,
	body string,
) error {
	var lastErr error

	for attempt := 1; attempt <= constant.DefaultMaxRetry; attempt++ {

		lastErr = c.send(
			toName,
			toEmail,
			subject,
			body,
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
	toName string,
	toEmail string,
	subject string,
	body string,
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
				TextPart: body,
				HTMLPart: body,
			},
		},
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.config.BaseURL+"/v3.1/send",
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

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf(
			"mailjet returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}

func (c *Client) SendWelcomeEmail(
	name string,
	email string,
) error {

	subject := "Welcome to Urbioxe!"

	body := fmt.Sprintf(`
		<h2>Welcome, %s 👋</h2>

		<p>
			Thank you for registering at <b>Urbioxe</b>.
		</p>

		<p>
			Your account has been created successfully.
		</p>

		<p>
			You can now report city issues,
			track your reports,
			and stay updated with your District or City information.
		</p>

		<br>

		<p>
			Regards,<br>
			Urbioxe Team
		</p>
	`, name)

	return c.SendEmail(
		name,
		email,
		subject,
		body,
	)
}

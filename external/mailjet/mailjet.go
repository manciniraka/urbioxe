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

func (c *Client) SendReportCreatedEmail(
	toEmail string,
	toName string,
	reportID int64,
	reportTitle string,
	categoryName string,
	addressLandmark string,
	createdAt string,
) error {
	subject := fmt.Sprintf("Laporan #%d Berhasil Dibuat: %s", reportID, reportTitle)

	textBody := fmt.Sprintf(
		"Halo %s,\n\nTerima kasih telah menyampaikan laporan/aduan Anda.\nLaporan Anda dengan ID #%d (%s) pada kategori %s telah berhasil dibuat pada %s.\nLokasi/Patokan: %s.\n\nLaporan Anda saat ini sedang dalam status PENDING dan akan ditinjau oleh tim admin.",
		toName, reportID, reportTitle, categoryName, createdAt, addressLandmark,
	)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
			<div style="background-color: #27ae60; color: white; padding: 15px; border-radius: 6px 6px 0 0; text-align: center;">
				<h2 style="margin: 0;">Laporan Berhasil Dibuat!</h2>
			</div>
			
			<div style="padding: 20px 0;">
				<p>Halo <strong>%s</strong>,</p>
				<p>Terima kasih telah berkontribusi aktif melaporkan masalah di sekitarmu. Laporanmu telah berhasil terdaftar di sistem kami dan akan segera ditinjau oleh dinas terkait.</p>
				
				<div style="background-color: #f8f9fa; padding: 15px; border-left: 4px solid #27ae60; margin: 20px 0;">
					<h3 style="margin-top: 0; color: #2c3e50;">Detail Laporan #%d</h3>
					<table style="width: 100%%; border-collapse: collapse; font-size: 14px;">
						<tr>
							<td style="padding: 6px 0; font-weight: bold; width: 35%%;">Judul Laporan:</td>
							<td style="padding: 6px 0;">%s</td>
						</tr>
						<tr>
							<td style="padding: 6px 0; font-weight: bold;">Kategori:</td>
							<td style="padding: 6px 0;">%s</td>
						</tr>
						<tr>
							<td style="padding: 6px 0; font-weight: bold;">Tanggal Dibuat:</td>
							<td style="padding: 6px 0;">%s</td>
						</tr>
						<tr>
							<td style="padding: 6px 0; font-weight: bold;">Lokasi / Patokan:</td>
							<td style="padding: 6px 0;">%s</td>
						</tr>
						<tr>
							<td style="padding: 6px 0; font-weight: bold;">Status Awal:</td>
							<td style="padding: 6px 0;"><span style="background-color: #f39c12; color: white; padding: 2px 8px; border-radius: 4px; font-size: 12px; font-weight: bold;">PENDING</span></td>
						</tr>
					</table>
				</div>

				<p style="font-size: 13px; color: #555;">Kamu akan menerima email notifikasi secara otomatis setiap kali ada pembaruan status atau tindak lanjut dari petugas lapangan terhadap laporan ini.</p>
			</div>

			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;" />
			<small style="color: #7f8c8d; text-align: center; display: block;">Email ini dikirim otomatis oleh Sistem Layanan Pengaduan Masyarakat. Mohon untuk tidak membalas email ini.</small>
		</div>
	`, toName, reportID, reportTitle, categoryName, createdAt, addressLandmark)

	return c.SendEmail(toEmail, toName, subject, textBody, htmlBody)
}

func (c *Client) SendReportStatusEmail(toEmail, toName, reportTitle string, reportID int64, status, notes string) error {
	subject := fmt.Sprintf("Update Status Laporan #%d: %s", reportID, reportTitle)

	textBody := fmt.Sprintf(
		"Halo %s,\n\nLaporan Anda #%d (%s) telah diperbarui menjadi status: %s.\nCatatan: %s",
		toName, reportID, reportTitle, status, notes,
	)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
			<h2 style="color: #2c3e50;">Pemberitahuan Status Laporan</h2>
			<p>Halo <strong>%s</strong>,</p>
			<p>Ada pembaruan status untuk laporan yang kamu ajukan:</p>
			
			<table style="width: 100%%; margin: 20px 0; border-collapse: collapse;">
				<tr>
					<td style="padding: 8px; font-weight: bold; width: 30%%;">ID Laporan:</td>
					<td style="padding: 8px;">#%d</td>
				</tr>
				<tr>
					<td style="padding: 8px; font-weight: bold;">Judul:</td>
					<td style="padding: 8px;">%s</td>
				</tr>
				<tr>
					<td style="padding: 8px; font-weight: bold;">Status Terbaru:</td>
					<td style="padding: 8px; color: #2980b9; font-weight: bold; text-transform: uppercase;">%s</td>
				</tr>
				<tr>
					<td style="padding: 8px; font-weight: bold;">Catatan / Respon:</td>
					<td style="padding: 8px;">%s</td>
				</tr>
			</table>

			<p>Terima kasih telah berkontribusi dalam menjaga lingkungan dan pelayanan publik!</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;" />
			<small style="color: #7f8c8d;">Email ini dikirim otomatis oleh Sistem Layanan Pengaduan Masyarakat.</small>
		</div>
	`, toName, reportID, reportTitle, status, notes)

	return c.SendEmail(toEmail, toName, subject, textBody, htmlBody)
}

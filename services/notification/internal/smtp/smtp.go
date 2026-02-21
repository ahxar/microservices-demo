package smtp

import (
	"fmt"
	"net/smtp"
)

type SMTPClient struct {
	host string
	port string
	from string
}

func NewSMTPClient(host, port, from string) *SMTPClient {
	return &SMTPClient{
		host: host,
		port: port,
		from: from,
	}
}

func (c *SMTPClient) SendEmail(to, subject, htmlBody string) error {
	// Construct email message
	headers := make(map[string]string)
	headers["From"] = c.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%s", c.host, c.port)

	// For Mailhog (development), we don't need authentication
	err := smtp.SendMail(
		addr,
		nil, // No auth for Mailhog
		c.from,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

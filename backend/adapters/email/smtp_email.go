package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
)

// SMTPEmailSender sends emails via an SMTP server.
type SMTPEmailSender struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

func NewSMTPEmailSender(host string, port int, user, pass, from string) *SMTPEmailSender {
	if port == 0 {
		port = 587
	}
	return &SMTPEmailSender{Host: host, Port: port, User: user, Pass: pass, From: from}
}

// SendResetEmail sends a simple plaintext email containing the reset URL.
func (s *SMTPEmailSender) SendResetEmail(ctx context.Context, toEmail string, resetURL string) error {
	if s.Host == "" || s.Port == 0 {
		return fmt.Errorf("smtp not configured")
	}

	var addr string
	if strings.Contains(s.Host, ":") {
		// IPv6 address, enclose in brackets
		addr = fmt.Sprintf("[%s]:%d", s.Host, s.Port)
	} else {
		addr = fmt.Sprintf("%s:%d", s.Host, s.Port)
	}

	// warn if password not provided
	if s.User != "" && s.Pass == "" {
		log.Printf("[SMTPEmail] Warning: SMTP user provided but SMTP_PASS is empty; authentication will likely fail")
	}

	// Build message
	subject := "Password reset"
	body := fmt.Sprintf("You requested a password reset. Click the link below to reset your password:\n\n%s\n\nIf you didn't request this, you can ignore this email.", resetURL)
	// add Reply-To so replies go to the user (optional)
	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", s.From),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Reply-To: %s", toEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.User, s.Pass, s.Host)

	log.Printf("[SMTPEmail] Attempting to send to=%s via %s:%d (user=%t)", toEmail, s.Host, s.Port, s.User != "")
	// Establish connection. If port is 465, use TLS directly (SMTPS). Otherwise use plain TCP and attempt STARTTLS if available.
	var err error
	var c *smtp.Client
	if s.Port == 465 {
		var tlsConn net.Conn
		tlsConn, err = tls.Dial("tcp", addr, &tls.Config{ServerName: s.Host})
		if err != nil {
			return fmt.Errorf("failed to dial tls smtp: %w", err)
		}
		c, err = smtp.NewClient(tlsConn, s.Host)
		if err != nil {
			return fmt.Errorf("failed to create smtp client over tls: %w", err)
		}
	} else {
		var conn net.Conn
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to dial smtp: %w", err)
		}
		c, err = smtp.NewClient(conn, s.Host)
		if err != nil {
			return fmt.Errorf("failed to create smtp client: %w", err)
		}

		// If server supports STARTTLS, upgrade the connection
		if ok, _ := c.Extension("STARTTLS"); ok {
			tlsconfig := &tls.Config{ServerName: s.Host}
			if err = c.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("starttls failed: %w", err)
			}
		}
	}
	defer c.Close()

	// Perform auth only if credentials were provided
	if s.User != "" {
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
	}

	if err = c.Mail(s.From); err != nil {
		return fmt.Errorf("smtp mail from failed: %w", err)
	}
	if err = c.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp rcpt to failed: %w", err)
	}

	wc, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data failed: %w", err)
	}
	_, err = wc.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("smtp write failed: %w", err)
	}
	if err = wc.Close(); err != nil {
		return fmt.Errorf("smtp close failed: %w", err)
	}

	if err := c.Quit(); err != nil {
		return err
	}
	log.Printf("[SMTPEmail] Sent to=%s via %s:%d", toEmail, s.Host, s.Port)
	return nil
}

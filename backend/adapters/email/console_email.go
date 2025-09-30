package email

import (
	"context"
	"log"
)

// ConsoleEmailSender prints reset links to the application log (development only)
type ConsoleEmailSender struct{}

func NewConsoleEmailSender() *ConsoleEmailSender { return &ConsoleEmailSender{} }

// SendResetEmail implements the EmailSender interface used by the usecase
func (s *ConsoleEmailSender) SendResetEmail(ctx context.Context, toEmail string, resetURL string) error {
	log.Printf("[ConsoleEmail] To=%s ResetURL=%s", toEmail, resetURL)
	return nil
}

package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"whatsapp-parser/internal/domain"
	"whatsapp-parser/pkg/selenium"
)

type sessionUseCase struct {
	repo           domain.SessionRepository
	whatsappClient *selenium.WhatsAppClient
}

// NewSessionUseCase creates a new session use case
func NewSessionUseCase(repo domain.SessionRepository) (domain.SessionUseCase, error) {
	client, err := selenium.NewWhatsAppClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp client: %v", err)
	}

	return &sessionUseCase{
		repo:           repo,
		whatsappClient: client,
	}, nil
}

func (u *sessionUseCase) CreateSession() (*domain.Session, string, error) {
	// Get QR code
	qrCode, err := u.whatsappClient.GetQRCode(context.Background())
	if err != nil {
		return nil, "", fmt.Errorf("failed to get QR code: %v", err)
	}

	// Create new session
	session := &domain.Session{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save session
	if err := u.repo.Save(session); err != nil {
		return nil, "", fmt.Errorf("failed to save session: %v", err)
	}

	return session, qrCode, nil
}

func (u *sessionUseCase) RestoreSession(id string) error {
	// Get session from repository
	session, err := u.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	// Convert session data to bytes
	sessionData, err := u.whatsappClient.GetSessionData()
	if err != nil {
		return fmt.Errorf("failed to get session data: %v", err)
	}

	// Restore session in WhatsApp client
	if err := u.whatsappClient.RestoreSession(sessionData); err != nil {
		return fmt.Errorf("failed to restore session: %v", err)
	}

	return nil
}

func (u *sessionUseCase) SendMessage(sessionID string, phoneNumber string, message string) error {
	// Verify session exists
	session, err := u.repo.GetByID(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %v", err)
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	// Send message
	if err := u.whatsappClient.SendMessage(phoneNumber, message); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
} 
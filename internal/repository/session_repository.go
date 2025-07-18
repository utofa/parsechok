package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"whatsapp-parser/internal/domain"
)

type sessionRepository struct {
	storagePath string
	mu          sync.RWMutex
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(storagePath string) (domain.SessionRepository, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %v", err)
	}

	return &sessionRepository{
		storagePath: storagePath,
	}, nil
}

func (r *sessionRepository) Save(session *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}

	filePath := filepath.Join(r.storagePath, session.ID+".json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %v", err)
	}

	return nil
}

func (r *sessionRepository) GetByID(id string) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath := filepath.Join(r.storagePath, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read session file: %v", err)
	}

	var session domain.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %v", err)
	}

	return &session, nil
}

func (r *sessionRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.storagePath, id+".json")
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete session file: %v", err)
	}

	return nil
} 
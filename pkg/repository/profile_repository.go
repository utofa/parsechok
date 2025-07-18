package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"whatsapp-parser/pkg/domain"
)

// ProfileRepository defines the interface for WhatsApp profile operations
type ProfileRepository interface {
	Get(id int) (*domain.WhatsAppProfile, error)
	Save(profile *domain.WhatsAppProfile) error
	Validate(path string) (bool, string)
}

// FileProfileRepository implements ProfileRepository using file system storage
type FileProfileRepository struct {
	basePath string
}

// NewFileProfileRepository creates a new FileProfileRepository
func NewFileProfileRepository(basePath string) *FileProfileRepository {
	return &FileProfileRepository{
		basePath: basePath,
	}
}

// Get retrieves a WhatsApp profile by ID
func (r *FileProfileRepository) Get(id int) (*domain.WhatsAppProfile, error) {
	profilePath := filepath.Join(r.basePath, fmt.Sprintf("profile_%d", id))
	
	// Check if profile directory exists
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("profile not found: %d", id)
	}

	return &domain.WhatsAppProfile{
		ID:      id,
		Path:    profilePath,
		IsValid: true, // We'll assume it's valid if it exists
	}, nil
}

// Save saves a WhatsApp profile
func (r *FileProfileRepository) Save(profile *domain.WhatsAppProfile) error {
	profilePath := filepath.Join(r.basePath, fmt.Sprintf("profile_%d", profile.ID))
	
	// Create profile directory if it doesn't exist
	if err := os.MkdirAll(profilePath, 0755); err != nil {
		return fmt.Errorf("failed to create profile directory: %v", err)
	}

	return nil
}

// Validate checks if a profile path is valid
func (r *FileProfileRepository) Validate(path string) (bool, string) {
	// Check if directory exists and is accessible
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, "Profile directory does not exist"
		}
		return false, fmt.Sprintf("Failed to access profile directory: %v", err)
	}

	return true, ""
} 
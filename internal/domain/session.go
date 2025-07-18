package domain

import "time"

// Session represents a WhatsApp Web session
type Session struct {
	ID        string    `json:"id"`
	Cookies   []Cookie  `json:"cookies"`
	Storage   []Storage `json:"storage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Cookie represents browser cookie
type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	Secure   bool      `json:"secure"`
	HttpOnly bool      `json:"http_only"`
}

// Storage represents local storage item
type Storage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SessionRepository interface for session persistence
type SessionRepository interface {
	Save(session *Session) error
	GetByID(id string) (*Session, error)
	Delete(id string) error
}

// SessionUseCase interface for session business logic
type SessionUseCase interface {
	CreateSession() (*Session, string, error) // Returns session, QR code URL, and error
	RestoreSession(id string) error
	SendMessage(sessionID string, phoneNumber string, message string) error
} 
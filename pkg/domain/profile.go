package domain

// WhatsAppProfile represents a WhatsApp profile configuration
type WhatsAppProfile struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	IsValid  bool   `json:"is_valid"`
}

// MarkValid marks the profile as valid
func (p *WhatsAppProfile) MarkValid() {
	p.IsValid = true
} 
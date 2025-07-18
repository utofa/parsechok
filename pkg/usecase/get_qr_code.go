package usecase

import (
	"fmt"
	//"whatsapp-parser/pkg/domain"
	"whatsapp-parser/pkg/repository"
	"whatsapp-parser/pkg/selenium"
)

// GetQRCodeUseCase handles the business logic for getting WhatsApp QR codes
type GetQRCodeUseCase struct {
	profileRepository repository.ProfileRepository
	whatsappClient    *selenium.WhatsAppClient
}

// NewGetQRCodeUseCase creates a new GetQRCodeUseCase
func NewGetQRCodeUseCase(
	profileRepository repository.ProfileRepository,
	whatsappClient *selenium.WhatsAppClient,
) *GetQRCodeUseCase {
	return &GetQRCodeUseCase{
		profileRepository: profileRepository,
		whatsappClient:    whatsappClient,
	}
}

// Result represents the result of getting a QR code
type Result struct {
	QRCode string
	Error  string
}

// Execute handles the QR code retrieval process
func (uc *GetQRCodeUseCase) Execute(sessionID string, profileID int) (*Result, error) {
	// Get profile
	profile, err := uc.profileRepository.Get(profileID)
	if err != nil {
		return &Result{Error: "Profile not found"}, nil
	}

	// Validate profile if not valid
	if !profile.IsValid {
		isValid, message := uc.profileRepository.Validate(profile.Path)
		if !isValid {
			return &Result{Error: message}, nil
		}
		profile.MarkValid()
		if err := uc.profileRepository.Save(profile); err != nil {
			return nil, fmt.Errorf("failed to save profile: %v", err)
		}
	}

	// Get QR code
	qrCode, err := uc.whatsappClient.GetQRCode(sessionID)
	if err != nil {
		if err.Error() == "Already authorized" {
			return &Result{Error: "Already authorized"}, nil
		}
		return nil, fmt.Errorf("failed to get QR code: %v", err)
	}

	return &Result{
		QRCode: qrCode,
	}, nil
}

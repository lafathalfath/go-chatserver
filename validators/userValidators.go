package validators

import (
	"errors"
	"github.com/lafathalfath/go-chatserver/graph/models"

	"gorm.io/gorm"
)

func ValidateNewUser(db *gorm.DB, input models.NewUser) error {
	var errorMessages []error

	if input.Name == "" {
		errorMessages = append(errorMessages, errors.New("User name cannot be null"))
	}
	if input.Email == "" {
		errorMessages = append(errorMessages, errors.New("User email cannot be null"))
	}
	if input.Password == "" {
		errorMessages = append(errorMessages, errors.New("User password cannot be null"))
	}
	if input.PasswordConfirmation == "" {
		errorMessages = append(errorMessages, errors.New("User password not confirmed"))
	}
	if input.PasswordConfirmation != input.Password {
		errorMessages = append(errorMessages, errors.New("User password confirmation incorrect"))
	}
	var prevUserEmail models.User
	if err := db.Select("email").First(&prevUserEmail, "email = ?", input.Email).Error; err == nil {
		errorMessages = append(errorMessages, errors.New("User already exist"))
	}
	return errors.Join(errorMessages...)
}

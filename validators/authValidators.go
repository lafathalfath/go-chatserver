package validators

import (
	"errors"
	"github.com/lafathalfath/go-chatserver/graph/models"
	"github.com/lafathalfath/go-chatserver/helpers"

	"gorm.io/gorm"
)

func ValidateLogin(db *gorm.DB, input models.Login) (models.User, error) {
	var errs []error
	var user models.User
	if err := db.Select("email").First(&user, "email = ?", input.Email).Error; err != nil {
		errs = append(errs, errors.New("You didn't register yet"))
	} else {
		if input.Email == "" {
			errs = append(errs, errors.New("Email cannot be null"))
		}
		if input.Password == "" {
			errs = append(errs, errors.New("Password cannot be null"))
		} else {
			db.Select("id", "password").First(&user, "email = ?", input.Email)
			if err := helpers.ComparePassword(input.Password, user.Password); err != nil {
				errs = append(errs, errors.New("Password incorrect"))
			}
		}
	}
	return user, errors.Join(errs...)
}

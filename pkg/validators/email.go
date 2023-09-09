package validators

import "github.com/asaskevich/govalidator"

// EmailValidation validate email format
func (v *Validation) EmailValidation(email string) bool {
	return govalidator.IsEmail(email)
}

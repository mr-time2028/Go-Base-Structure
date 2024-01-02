package validators

import "github.com/asaskevich/govalidator"

// EmailValidation validate email format
func (v *Validation) EmailValidation(email string) {
	if !govalidator.IsEmail(email) {
		v.Errors.Add("email", "invalid email address")
	}
}

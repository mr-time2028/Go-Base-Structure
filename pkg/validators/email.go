package validators

import "github.com/asaskevich/govalidator"

func (v *Validation) EmailValidation(email string) bool {
	return govalidator.IsEmail(email)
}

package validators

import "github.com/asaskevich/govalidator"

func EmailValidation(email string) bool {
	return govalidator.IsEmail(email)
}

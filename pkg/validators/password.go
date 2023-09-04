package validators

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

type PassErrors map[string][]string

func newPassErrors() PassErrors {
	return make(map[string][]string)
}

func (p PassErrors) AddError(msg string) {
	p["password"] = append(p["password"], msg)
}

func PasswordMatchesValidation(hashedDBPassword, ClientPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedDBPassword), []byte(ClientPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, bcrypt.ErrMismatchedHashAndPassword
		default:
			return false, err
		}
	}

	return true, nil
}

func PasswordCharacterValidation(password string) PassErrors {
	passErrors := newPassErrors()

	hasUppercase := false
	hasLowercase := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		} else if unicode.IsLower(char) {
			hasLowercase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if len(password) < 8 {
		passErrors.AddError("password must be a minimum of 8 characters in length")
	}
	if len(password) > 30 {
		passErrors.AddError("password must be a maximum of 30 characters in length")
	}
	if !hasUppercase {
		passErrors.AddError("password must contain at least one uppercase letter")
	}
	if !hasLowercase {
		passErrors.AddError("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		passErrors.AddError("password must contain at least one digit")
	}

	return passErrors
}

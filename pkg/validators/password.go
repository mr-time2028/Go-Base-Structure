package validators

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode"
)

// PasswordMatchesValidation check if two password are match, password from db (registered) and password from client
func (v *Validation) PasswordMatchesValidation(hashedDBPassword, ClientPassword string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedDBPassword), []byte(ClientPassword))
	if err != nil {
		v.Errors.Add("password", err.Error())
		v.Errors.Code = http.StatusUnauthorized
	}
}

// PasswordCharacterValidation check if password characters are valid (min length, max length and etc.)
func (v *Validation) PasswordCharacterValidation(password string) {
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
		v.Errors.Add("password", "this field must be a minimum length of 8 characters")
	}
	if len(password) > 30 {
		v.Errors.Add("password", "this field must be a maximum length of 30 characters")
	}
	if !hasUppercase {
		v.Errors.Add("password", "this field must contain at least one uppercase letter")
	}
	if !hasLowercase {
		v.Errors.Add("password", "this field must contain at least one lowercase letter")
	}
	if !hasDigit {
		v.Errors.Add("password", "this field must contain at least one digit")
	}
}

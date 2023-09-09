package validators

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordMatchesValidation(t *testing.T) {
	// case 1: passwords are the same
	hashedDBPassword, _ := bcrypt.GenerateFromPassword([]byte("MyPassword1234"), bcrypt.DefaultCost)
	ClientPassword := "MyPassword1234"

	validator := New()
	isMatch := validator.PasswordMatchesValidation(string(hashedDBPassword), ClientPassword)

	if !isMatch {
		t.Errorf("the passwords should match, but do not")
	}
	if !validator.Valid() {
		t.Errorf("unexpected error: %s", validator.Errors.Get("password"))
	}

	// case 2: passwords are not the same
	hashedDBPassword, _ = bcrypt.GenerateFromPassword([]byte("MyPassword1234"), bcrypt.DefaultCost)
	ClientPassword = "MyPass"

	validator = New()
	isMatch = validator.PasswordMatchesValidation(string(hashedDBPassword), ClientPassword)

	if isMatch {
		t.Errorf("the passwords should not be match, but they do")
	}
	errMsg := validator.Errors.Get("password")
	if errMsg != bcrypt.ErrMismatchedHashAndPassword.Error() {
		t.Errorf("unexpected error, should get: %s, but got: %s", bcrypt.ErrMismatchedHashAndPassword.Error(), errMsg)
	}
}

func TestPasswordCharactersValidation(t *testing.T) {
	var testCases = []struct {
		name           string
		password       string
		expectedErr    bool
		expectedErrMsg string
	}{
		{"valid password", "Password1234", false, ""},
		{"invalid min length", "Pas2", true, "min length"},
		{"invalid max length", "Pasdafjhqwjhkjghdaugioruetlkjkllkjsgq2", true, "max length"},
		{"no uppercase letter", "abcdefghi158", true, "no uppercase letter"},
		{"no lowercase letter", "ABCDEFGHI158", true, "no lowercase letter"},
		{"no digit", "ADGWatcew", true, "no digit"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator := New()
			validator.PasswordCharacterValidation(tc.password)

			err := validator.Errors.Get("password")

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				t.Errorf("unexpected error: %v", err)
			} else if tc.expectedErr && validator.Valid() {
				t.Errorf("expected a %s error, but got nil", tc.expectedErrMsg)
			}
		})
	}
}

package validators

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordMatchesValidation(t *testing.T) {
	// case 1: passwords are the same
	hashedDBPassword, _ := bcrypt.GenerateFromPassword([]byte("MyPassword1234"), bcrypt.DefaultCost)
	ClientPassword := "MyPassword1234"

	isMatch, err := PasswordMatchesValidation(string(hashedDBPassword), ClientPassword)
	if isMatch == false {
		t.Errorf("the passwords should match, but don't and got error: %s", err)
	}

	// case 2: passwords are not the same
	hashedDBPassword, _ = bcrypt.GenerateFromPassword([]byte("MyPassword1234"), bcrypt.DefaultCost)
	ClientPassword = "MyPass"

	isMatch, err = PasswordMatchesValidation(string(hashedDBPassword), ClientPassword)
	if isMatch {
		t.Errorf("the passwords should not be match, but they do: %s", err)
	}
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Errorf("unexpected error, should get: %s, but got: %s", bcrypt.ErrMismatchedHashAndPassword.Error(), err)
	}
}

func TestPasswordCharactersValidation(t *testing.T) {
	var testCases = []struct {
		name           string
		password       string
		expectedErr    bool
		expectedErrStr string
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
			err := PasswordCharacterValidation(tc.password)
			if tc.expectedErr && len(err) == 0 || !tc.expectedErr && len(err) > 0 {
				t.Errorf("%s: unexpected error: %v", tc.name, err)
			} else if tc.expectedErr && len(err) == 0 {
				t.Errorf("%s: expected a %s error, but got nil", tc.name, tc.expectedErrStr)
			}
		})
	}
}

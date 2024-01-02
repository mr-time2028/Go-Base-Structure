package validators

import "testing"


func TestValidation_EmailValidation(t *testing.T) {
	theTests := []struct {
		name        string
		email		string
		expectedErr bool
	} {
		{
			"valid email",
			"John@valid.com",
			false,
		},
		{
			"invalid email",
			"John.com",
			true,
		},
	}

	for _, tc := range theTests {
		t.Run(tc.name, func(t *testing.T) {
			validator := New()
			validator.EmailValidation(tc.email)

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				err := validator.Errors.Get("email")
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}

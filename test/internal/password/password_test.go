package password_test

import (
	"testing"
	"uacl/internal/password"

	"github.com/stretchr/testify/assert"
)

func TestCreatePassword(t *testing.T) {
	password := password.CreatePassword("Test123")

	assert.NotNil(t, password)
}

func TestVerifyPassword(t *testing.T) {
	type testCase struct {
		useGeneratedPassword bool
		want                 bool
	}

	tests := []testCase{
		{
			useGeneratedPassword: true,
			want:                 true,
		},
		{
			useGeneratedPassword: false,
			want:                 false,
		},
	}

	for _, tc := range tests {
		generatedPassword := password.CreatePassword("Test123")

		var result bool

		if tc.useGeneratedPassword {
			result = password.ValidatePassword("Test123", generatedPassword)
		} else {
			result = password.ValidatePassword("Test12", generatedPassword)
		}

		assert.Equal(t, tc.want, result)
	}
}

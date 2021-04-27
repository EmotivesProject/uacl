package password_test

import (
	"testing"
	"uacl/internal/password"

	"github.com/stretchr/testify/assert"
)

func TestPasswordCanCreate(t *testing.T) {
	result := password.CreatePassword("test123")
	assert.NotNil(
		t,
		result,
	)
}

func TestVerifyPassword(t *testing.T) {
	type passwords struct {
		password     string
		testPassword string
		want         bool
	}

	tests := []passwords{
		{
			password:     "test123",
			testPassword: "test123",
			want:         true,
		},
		{
			password:     "test123",
			testPassword: "test124",
			want:         false,
		},
	}

	for _, tc := range tests {
		hashedPassword := password.CreatePassword(tc.password)
		result := password.ValidatePassword(tc.testPassword, hashedPassword)
		assert.Equal(t, tc.want, result)
	}
}

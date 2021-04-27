package auth_test

import (
	"testing"
	"uacl/internal/auth"
	"uacl/model"

	"github.com/stretchr/testify/assert"
)

var testUser = model.User{
	ID:       1,
	Name:     "Jane",
	Username: "TestAcc",
}

func TestAuthCanCreate(t *testing.T) {
	_, err := auth.CreateToken(testUser)
	assert.NotNil(
		t,
		err,
	)
}

func TestVerifyAuth(t *testing.T) {
	type authTest struct {
		use bool
	}

	tests := []authTest{
		{
			use: false,
		},
	}

	for _, tc := range tests {
		token, _ := auth.CreateToken(testUser)
		if tc.use {
			_, err := auth.Validate(token)
			assert.Nil(t, err)
		} else {
			_, err := auth.Validate("token")
			assert.NotNil(t, err)
		}
	}
}

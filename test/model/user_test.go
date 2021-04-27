package model_test

import (
	"testing"
	"uacl/model"
	"uacl/test"

	"github.com/stretchr/testify/assert"
)

var testUser = model.User{}

func setupTestCase() {
	testUser.Name = "Jane"
	testUser.Username = "JaneTT"
	testUser.Password = "Test123"
	testUser.Secret = "qutCreate"
}

func TestNameValidLength(t *testing.T) {
	type nameLength struct {
		length int
		want   bool
	}

	tests := []nameLength{
		{
			length: 2,
			want:   false,
		},
		{
			length: 5,
			want:   true,
		},
		{
			length: 101,
			want:   false,
		},
	}

	for _, tc := range tests {
		setupTestCase()

		newName := test.CreateStringAtLength(tc.length)
		testUser.Name = newName

		_, err := testUser.ValidateCreate()
		result := err == nil

		assert.Equal(t, tc.want, result)
	}
}

func TestUsernameValidLength(t *testing.T) {
	type usernameLength struct {
		length int
		want   bool
	}

	tests := []usernameLength{
		{
			length: 2,
			want:   false,
		},
		{
			length: 5,
			want:   true,
		},
		{
			length: 101,
			want:   false,
		},
	}

	for _, tc := range tests {
		setupTestCase()

		newUsername := test.CreateStringAtLength(tc.length)
		testUser.Username = newUsername

		_, err := testUser.ValidateCreate()
		result := err == nil

		assert.Equal(t, tc.want, result)
	}
}

func TestPasswordValidLength(t *testing.T) {
	type passwordLength struct {
		length int
		want   bool
	}

	tests := []passwordLength{
		{
			length: 5,
			want:   false,
		},
		{
			length: 7,
			want:   true,
		},
		{
			length: 101,
			want:   false,
		},
	}

	for _, tc := range tests {
		setupTestCase()

		newPassword := test.CreateStringAtLength(tc.length)
		testUser.Password = newPassword

		_, err := testUser.ValidateCreate()
		result := err == nil

		assert.Equal(t, tc.want, result)
	}
}

func TestInvalidCharacters(t *testing.T) {
	type invalidCharacter struct {
		username string
		want     bool
	}

	tests := []invalidCharacter{
		{
			username: "test",
			want:     true,
		},
		{
			username: "test@#@!@213",
			want:     false,
		},
	}

	for _, tc := range tests {
		setupTestCase()

		testUser.Username = tc.username

		_, err := testUser.ValidateCreate()
		result := err == nil

		assert.Equal(t, tc.want, result)
	}
}

func TestRouterInitializes(t *testing.T) {
	setupTestCase()

	testUser.Secret = "incorrectSecret"

	_, err := testUser.ValidateCreate()
	result := err == nil

	assert.False(t, result)
}

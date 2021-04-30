// +build integration

package auth_test

import (
	"testing"
	"uacl/internal/auth"
	"uacl/model"
	"uacl/test"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateToken(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:     "test_acc",
		Username: "test_acc",
		Password: "test123",
	}

	_, err := auth.CreateToken(user)
	assert.Nil(t, err)

	test.TearDownIntegrationTest()
}

// +build integration

package db_test

import (
	"testing"
	"uacl/internal/db"
	"uacl/model"
	"uacl/test"

	"github.com/stretchr/testify/assert"
)

func TestCanFindUserByUsername(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:     "test_acc",
		Username: "test_acc_3",
		Password: "test123",
	}
	err := db.CreateNewUser(&user)
	assert.Nil(t, err)

	_, err = db.FindByUsername(user.Username)
	assert.Nil(t, err)

	test.TearDownIntegrationTest()
}

func TestCanNotFindUserByUsername(t *testing.T) {
	test.SetUpIntegrationTest()

	_, err := db.FindByUsername("fake_user")
	assert.NotNil(t, err)

	test.TearDownIntegrationTest()
}

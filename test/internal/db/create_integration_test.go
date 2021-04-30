// +build integration

package db_test

import (
	"testing"
	"uacl/internal/db"
	"uacl/model"
	"uacl/test"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateUser(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:     "test_acc",
		Username: "test_acc",
		Password: "test123",
	}
	err := db.CreateNewUser(&user)
	assert.Nil(t, err)

	test.TearDownIntegrationTest()
}

func TestCanNotCreateSameUser(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:     "test_acc",
		Username: "test_acc_2",
		Password: "test123",
	}
	err := db.CreateNewUser(&user)
	assert.Nil(t, err)
	err = db.CreateNewUser(&user)
	assert.NotNil(t, err)

	test.TearDownIntegrationTest()
}

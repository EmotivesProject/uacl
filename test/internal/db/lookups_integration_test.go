// +build integration

package db_test

import (
	"context"
	"testing"
	"uacl/internal/db"
	"uacl/model"
	"uacl/test"

	"github.com/stretchr/testify/assert"
)

func TestCanFindUserByUsername(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:      "test_acc",
		Username:  "test_acc_3",
		Password:  "test123",
		UserGroup: "qut",
	}
	err := db.CreateNewUser(context.Background(), &user)
	assert.Nil(t, err)

	_, err = db.FindByUsername(context.Background(), user.Username)
	assert.Nil(t, err)

	test.TearDownIntegrationTest()
}

func TestCanNotFindUserByUsername(t *testing.T) {
	test.SetUpIntegrationTest()

	_, err := db.FindByUsername(context.Background(), "fake_user")
	assert.NotNil(t, err)

	test.TearDownIntegrationTest()
}

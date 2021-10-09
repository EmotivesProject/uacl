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

func TestCanCreateUser(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:      "test_acc",
		Username:  "test_acc",
		Password:  "test123",
		UserGroup: "qut",
	}
	err := db.CreateNewUser(context.Background(), &user)
	assert.Nil(t, err)

	test.TearDownIntegrationTest()
}

func TestCanNotCreateSameUser(t *testing.T) {
	test.SetUpIntegrationTest()

	user := model.User{
		Name:      "test_acc",
		Username:  "test_acc_2",
		Password:  "test123",
		UserGroup: "qut",
	}
	err := db.CreateNewUser(context.Background(), &user)
	assert.Nil(t, err)
	err = db.CreateNewUser(context.Background(), &user)
	assert.NotNil(t, err)

	test.TearDownIntegrationTest()
}

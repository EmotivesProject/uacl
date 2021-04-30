// +build integration

package db_test

import (
	"testing"
	"uacl/internal/db"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/stretchr/testify/assert"
)

func TestRouterInitializes(t *testing.T) {
	logger.InitLogger("uacl")

	commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/uacl_db",
	})

	user := model.User{
		Name:     "test_acc",
		Username: "test_acc",
		Password: "test123",
	}

	err := db.CreateNewUser(&user)

	assert.Nil(t, err)
}

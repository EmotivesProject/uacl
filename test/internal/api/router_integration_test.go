// +build integration

package api_test

import (
	"net/http"
	"strings"
	"testing"
	"uacl/test"

	"github.com/stretchr/testify/assert"
)

func TestRouterHealthzHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	_, body := test.TestRequest(t, test.TS, "GET", "/healthz", nil)

	assert.EqualValues(t, "{\"result\":null,\"message\":[{\"message\":\"Health OK\"}]}", body)

	test.TearDownIntegrationTest()
}

func TestRouterUserHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	requestBody := strings.NewReader(
		"{\"username\": \"imtom15\", \"name\": \"imtom\", \"password\": \"test123\", \"secret\": \"qutCreate\" }",
	)

	r, _ := test.TestRequest(t, test.TS, "POST", "/user", requestBody)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

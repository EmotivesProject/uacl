// +build integration

package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"uacl/test"

	commonTest "github.com/EmotivesProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRouterHealthzHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	r, _, _ := commonTest.TestRequest(t, test.TS, "GET", "/healthz", nil)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterUserHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	requestBody := strings.NewReader(
		"{\"username\": \"imtom124\", \"name\": \"imtom\", \"password\": \"test123\", \"secret\": \"test\", \"user_group\": \"qut\" }",
	)

	r, _, _ := commonTest.TestRequest(t, test.TS, "POST", "/user", requestBody)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterAuthHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	requestBody := strings.NewReader(
		"{\"username\": \"imtom125\", \"name\": \"imtom\", \"password\": \"test123\", \"secret\": \"test\", \"user_group\": \"qut\" }",
	)

	_, resp, _ := commonTest.TestRequest(t, test.TS, "POST", "/user", requestBody)

	req, _ := http.NewRequest("GET", test.TS.URL+"/authorize", nil)
	req.Header.Add("Authorization", "Bearer "+resp["token"].(string))

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterRefreshHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	requestBody := strings.NewReader(
		"{\"username\": \"imtom132\", \"name\": \"imtom\", \"password\": \"test123\", \"secret\": \"test\", \"user_group\": \"qut\" }",
	)
	_, resp, _ := commonTest.TestRequest(t, test.TS, "POST", "/user", requestBody)

	requestBody = strings.NewReader(
		fmt.Sprintf("{\"username\": \"%s\", \"refresh_token\": \"%s\" }", resp["username"], resp["refresh_token"]),
	)
	req, _ := http.NewRequest("POST", test.TS.URL+"/refresh", requestBody)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterLoginHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	requestBody := strings.NewReader(
		"{\"username\": \"imtom134\", \"name\": \"imtom\", \"password\": \"test123\", \"secret\": \"test\", \"user_group\": \"qut\" }",
	)
	commonTest.TestRequest(t, test.TS, "POST", "/user", requestBody)

	requestBody = strings.NewReader(
		fmt.Sprintf("{\"username\": \"%s\", \"password\": \"%s\" }", "imtom134", "test123"),
	)
	req, _ := http.NewRequest("POST", test.TS.URL+"/login", requestBody)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

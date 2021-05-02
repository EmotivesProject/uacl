package test

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"uacl/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/response"
)

var TS *httptest.Server

func CreateStringAtLength(length int) string {
	var str string
	for i := 0; i < length; i++ {
		str += "x"
	}

	return str
}

func SetUpIntegrationTest() {
	logger.InitLogger("uacl")

	commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/uacl_db",
	})

	// These are set to be correct for internal/* tests
	os.Setenv("PRIVATE_KEY", "./../../../jwt/private.key")
	os.Setenv("PUBLIC_KEY", "./../../../jwt/public.key")

	os.Setenv("SECRET", "qutCreate")

	router := api.CreateRouter()

	TS = httptest.NewServer(router)
}

func TearDownIntegrationTest() {
	con := commonPostgres.GetDatabase()
	con.Close(context.TODO())

	TS.Close()
}

func TestRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	body io.Reader,
) (
	*http.Response, map[string]interface{}, []response.Message,
) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)

		return nil, nil, nil
	}

	return CompleteTestRequest(t, req)
}

func CompleteTestRequest(t *testing.T, r *http.Request) (*http.Response, map[string]interface{}, []response.Message) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)

		return nil, nil, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)

		return nil, nil, nil
	}
	defer resp.Body.Close()

	var responseObj response.Response

	err = json.Unmarshal(respBody, &responseObj)
	if err != nil {
		return resp, nil, responseObj.Message
	}

	resultMap, ok := responseObj.Result.(map[string]interface{})

	if !ok {
		return resp, nil, nil
	}

	return resp, resultMap, responseObj.Message
}

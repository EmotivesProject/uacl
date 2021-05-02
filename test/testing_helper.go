package test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"uacl/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
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

	router := api.CreateRouter()

	TS = httptest.NewServer(router)
}

func TearDownIntegrationTest() {
	con := commonPostgres.GetDatabase()
	con.Close(context.TODO())

	TS.Close()
}

func TestRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)

		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)

		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)

		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

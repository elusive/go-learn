package service_test

import (
    "fmt"
    "io"
    "os"
    "testing"

    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
    "github.com/elusive/instrument-api/server/service"
    "github.com/elusive/instrument-api/server/svclog"
    "v.io/x/lib/vlog"
)

type testContext struct {
    router *gin.Engine
}

var context = &testContext{}

func TestMain(m *testing.M) {
    svclog.Start()
    code := 1

    if testing.Verbose() {
		// Echo all server output to os.Stderr.
		go io.Copy(os.Stderr, ms.Stdout)
		go io.Copy(os.Stderr, ms.Stderr)
	}

    context.router = service.Register()
    code = m.Run()

    teardown()
    os.Exit(code)
} 

func TestLoadMethod(t *testing.T) {
    var createResult
	code, err := execute("POST", "/instrument/v1/load-method", poolingBatchCreate, &createResult)
	assert.Equal(t, http.StatusOK, code)
	assert.Nil(t, err)
}

func teardown() {
	// add here any cleanup as needed
}

func execute(method string, url string, req interface{}, resultPtr interface{}) (int, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}
	vlog.Info("Sending request: ", method, url, string(reqBytes))
	request, _ := http.NewRequest(method, url, bytes.NewReader(reqBytes))
	request.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	context.router.ServeHTTP(resp, request)

	vlog.Info("Received response:", resp.Body.String())
	err = json.Unmarshal(resp.Body.Bytes(), resultPtr)

	return resp.Code, err
}
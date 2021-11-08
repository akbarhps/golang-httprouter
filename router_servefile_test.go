package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

//go:embed resources
var resources embed.FS

func TestRouterServeFile(t *testing.T) {
	router := httprouter.New()

	// create a sub dir destination so we don't need
	// to specify the file folder, i.e resources
	dir, _ := fs.Sub(resources, "resources")

	// *filepath is hardcoded in ServeFile
	// so its must be named *filepath
	router.ServeFiles("/files/*filepath", http.FS(dir))

	request := httptest.NewRequest(http.MethodGet, "/files/test.txt", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "test", string(responseBody))

	t.Log(string(responseBody))
}

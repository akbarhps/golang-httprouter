package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprintf(writer, "Hello HTTP Router")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200")
	assert.Equal(t, "Hello HTTP Router", string(responseBody), "response body should be 'Hello HTTP Router'")
	t.Log(string(responseBody))
}

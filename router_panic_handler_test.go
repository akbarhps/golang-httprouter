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

func TestRouterPanicHandler(t *testing.T) {
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		t.Logf("Panic : %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error\nError : %s", err)
	}
	router.GET("/panic", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		panic("Oops, something went wrong")
	})

	request := httptest.NewRequest("GET", "/panic", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "Internal Server Error\nError : Oops, something went wrong", string(responseBody))

	t.Log(string(responseBody))
}

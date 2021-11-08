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

func TestRouterParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/products/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		t.Log(params)
		fmt.Fprintf(writer, "You requested a product with id %s", params.ByName("id"))
	})

	productId := "1"
	request := httptest.NewRequest(http.MethodGet, "/products/"+productId, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200")
	assert.Equal(t, "You requested a product with id "+productId, string(responseBody), "response body should be 'You requested a product with id 1'")
	t.Log(string(responseBody))
}

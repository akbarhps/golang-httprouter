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

func TestRouterPatternNamedParameter(t *testing.T) {
	router := httprouter.New()
	router.GET("/products/:id/items/:itemId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		t.Log(params)

		id := params.ByName("id")
		itemId := params.ByName("itemId")
		fmt.Fprintf(writer, "product id: %s, item id: %s", id, itemId)
	})

	request := httptest.NewRequest("GET", "/products/1/items/2", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "product id: 1, item id: 2", string(responseBody))

	t.Log(string(responseBody))
}

func TestRouterPatternCatchAllParameter(t *testing.T) {
	router := httprouter.New()
	router.GET("/images/*image", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		t.Log(params)

		image := params.ByName("image")
		fmt.Fprintf(writer, "image: %s", image)
	})

	request := httptest.NewRequest("GET", "/images/resources/img.jpg", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "image: /resources/img.jpg", string(responseBody))

	t.Log(string(responseBody))
}

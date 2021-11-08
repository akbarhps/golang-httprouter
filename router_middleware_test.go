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

type Logger struct {
	Handler http.Handler
}

func (log *Logger) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Before handler execution")
	log.Handler.ServeHTTP(w, r)
	fmt.Println("After handler execution")
}

func TestRouterLogMiddleware(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello world!")
	})

	logger := &Logger{
		Handler: router,
	}

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	logger.ServerHTTP(response, request)
	result := response.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	resultBody, err := ioutil.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello world!", string(resultBody))

	t.Log(string(resultBody))
}

package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	router.GET("/ping", Ping)
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	router.ServeHTTP(recorder, req)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, "Pong", recorder.Body.String())
}

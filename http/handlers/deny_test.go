package handlers

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDenyHandler(t *testing.T) {
	assert := require.New(t)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/v1/test", bytes.NewBufferString("TestBody"))

	DenyHandler(recorder, request)

	assert.EqualValues("{\"error\":\"Method not allowed\"}\n", recorder.Body.String())
	assert.EqualValues(http.StatusMethodNotAllowed, recorder.Code)
}

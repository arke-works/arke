package handlers

import (
	"bytes"
	"github.com/pressly/chi"
	chim "github.com/pressly/chi/middleware"
	"github.com/stretchr/testify/require"
	"iris.arke.works/forum/http/helper"
	"iris.arke.works/forum/http/middleware"
	"iris.arke.works/forum/http/resources"
	"iris.arke.works/forum/snowflakes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostHandler(t *testing.T) {
	assert := require.New(t)
	fountain := &snowflakes.SequenceGenerator{Position: 0}
	helper.SetupTestLog()

	router := chi.NewRouter()

	router.Use(chim.Recoverer, chim.RedirectSlashes, middleware.FountainMiddleware(fountain))

	router.Route("/", MakeRouter)

	resources.SetupMockDatatype()

	m := resources.MockResource{
		SnowflakeField: 22,
		TextField:      "Stripped",
		OtherTextField: "Not Stripped",
		IntField:       -23434,
		SliceField:     []byte{0x22, 0x00, 0xFF, 0x44},
		TimeField:      time.Now(),
		OptTimeField:   nil,
	}

	dat, err := m.MarshalJSON()
	assert.NoError(err)
	if err != nil {
		return
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/mock", bytes.NewBuffer(dat))

	PostHandler(recorder, request)

	assert.EqualValues(http.StatusBadRequest, recorder.Code)
	assert.EqualValues("", recorder.Body.String())

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/notaresource", bytes.NewBufferString(""))

	router.ServeHTTP(recorder, request)

	assert.EqualValues("{\"error\":\"Endpoint not registered\"}\n", recorder.Body.String())
	assert.EqualValues(http.StatusNotFound, recorder.Code)

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/mock_nofactory", bytes.NewBufferString("{}"))

	router.ServeHTTP(recorder, request)

	assert.EqualValues("{\"error\":\"Endpoint not registered\"}\n", recorder.Body.String())
	assert.EqualValues(http.StatusNotFound, recorder.Code)

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/mock", bytes.NewBuffer(dat))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(`{"id_field":1,"text":"","other_text":"Not Stripped","int":-23434,"bytes":"IgD/RA=="}` + "\n", recorder.Body.String())
	assert.EqualValues(http.StatusOK, recorder.Code)
}

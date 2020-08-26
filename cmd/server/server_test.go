package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hectorgabucio/genity-hector/internal/data"
	"github.com/hectorgabucio/genity-hector/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const TITLE = "title"

func TestGetData(t *testing.T) {
	tests := []struct {
		method         string
		dataRepository *mocks.DataRepository
		titleParam     string
		statusCode     int
	}{
		{http.MethodPost, mockRepositoryGetDataOK(), "", 405},
		{http.MethodGet, mockRepositoryGetDataOK(), "", 400},
		{http.MethodGet, mockRepositoryGetDataNotFound(), TITLE, 404},
		{http.MethodGet, mockRepositoryGetDataKO(), TITLE, 500},
		{http.MethodGet, mockRepositoryGetDataOK(), TITLE, 200},
	}

	assert := assert.New(t)
	for _, tt := range tests {
		app := &app{dataRepository: tt.dataRepository}
		testHandler, rr, req := prepareSUTGetData(tt.method, t, app)
		req.URL.Path = req.URL.Path + tt.titleParam
		testHandler.ServeHTTP(rr, req)

		assert.Equal(tt.statusCode, rr.Code, "handler returned wrong status code: got %v want %v",
			rr.Code, tt.statusCode)
	}
}

func prepareSUTGetData(method string, t *testing.T, app *app) (http.Handler, *httptest.ResponseRecorder, *http.Request) {
	handler := http.HandlerFunc(app.getData)

	rr := httptest.NewRecorder()

	req, err := http.NewRequest(method, "/get-data/", nil)
	if err != nil {
		t.Fatal(err)
	}

	return handler, rr, req

}

func mockRepositoryGetDataOK() *mocks.DataRepository {
	mockRepo := &mocks.DataRepository{}
	retrieved := &data.Data{Title: TITLE}
	mockRepo.On("Get", mock.Anything).Return(retrieved, nil)
	return mockRepo
}

func mockRepositoryGetDataKO() *mocks.DataRepository {
	mockRepo := &mocks.DataRepository{}
	mockRepo.On("Get", mock.Anything).Return(nil, errors.New("error"))
	return mockRepo
}

func mockRepositoryGetDataNotFound() *mocks.DataRepository {
	mockRepo := &mocks.DataRepository{}
	mockRepo.On("Get", mock.Anything).Return(nil, nil)
	return mockRepo
}

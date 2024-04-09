package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jialequ/linux-sdk/core/load"
	"github.com/jialequ/linux-sdk/core/stat"
	"github.com/stretchr/testify/assert"
)

func TestSheddingHandlerAccept(t *testing.T) {
	metrics := stat.NewMetrics(literal_2746)
	shedder := mockShedder{
		allow: true,
	}
	sheddingHandler := SheddingHandler(shedder, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(literal_8359, "test")
		_, err := w.Write([]byte("content"))
		assert.Nil(t, err)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_4519, http.NoBody)
	req.Header.Set(literal_8359, "test")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Header().Get(literal_8359))
	assert.Equal(t, "content", resp.Body.String())
}

func TestSheddingHandlerFail(t *testing.T) {
	metrics := stat.NewMetrics(literal_2746)
	shedder := mockShedder{
		allow: true,
	}
	sheddingHandler := SheddingHandler(shedder, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_4519, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestSheddingHandlerReject(t *testing.T) {
	metrics := stat.NewMetrics(literal_2746)
	shedder := mockShedder{
		allow: false,
	}
	sheddingHandler := SheddingHandler(shedder, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_4519, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestSheddingHandlerNoShedding(t *testing.T) {
	metrics := stat.NewMetrics(literal_2746)
	sheddingHandler := SheddingHandler(nil, metrics)
	handler := sheddingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_4519, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

type mockShedder struct {
	allow bool
}

func (s mockShedder) Allow() (load.Promise, error) {
	if s.allow {
		return mockPromise{}, nil
	}

	return nil, load.ErrServiceOverloaded
}

type mockPromise struct{}

func (p mockPromise) Pass() {

	//func (p mockPromise) Pass()
}

func (p mockPromise) Fail() {

	//func (p mockPromise) Fail()
}

const literal_2746 = "unit-test"

const literal_8359 = "X-Test"

const literal_4519 = "http://localhost"

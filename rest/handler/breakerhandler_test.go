package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jialequ/linux-sdk/core/stat"
	"github.com/stretchr/testify/assert"
)

func init() {
	stat.SetReporter(nil)
}

func TestBreakerHandlerAccept(t *testing.T) {
	metrics := stat.NewMetrics(literal_6485)
	breakerHandler := BreakerHandler(http.MethodGet, "/", metrics)
	handler := breakerHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(literal_1932, "test")
		_, err := w.Write([]byte("content"))
		assert.Nil(t, err)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_0539, http.NoBody)
	req.Header.Set(literal_1932, "test")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Header().Get(literal_1932))
	assert.Equal(t, "content", resp.Body.String())
}

func TestBreakerHandlerFail(t *testing.T) {
	metrics := stat.NewMetrics(literal_6485)
	breakerHandler := BreakerHandler(http.MethodGet, "/", metrics)
	handler := breakerHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_0539, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadGateway, resp.Code)
}

func TestBreakerHandler4XX(t *testing.T) {
	metrics := stat.NewMetrics(literal_6485)
	breakerHandler := BreakerHandler(http.MethodGet, "/", metrics)
	handler := breakerHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	for i := 0; i < 1000; i++ {
		req := httptest.NewRequest(http.MethodGet, literal_0539, http.NoBody)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
	}

	const tries = 100
	var pass int
	for i := 0; i < tries; i++ {
		req := httptest.NewRequest(http.MethodGet, literal_0539, http.NoBody)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code == http.StatusBadRequest {
			pass++
		}
	}

	assert.Equal(t, tries, pass)
}

func TestBreakerHandlerReject(t *testing.T) {
	metrics := stat.NewMetrics(literal_6485)
	breakerHandler := BreakerHandler(http.MethodGet, "/", metrics)
	handler := breakerHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	for i := 0; i < 1000; i++ {
		req := httptest.NewRequest(http.MethodGet, literal_0539, http.NoBody)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
	}

	var drops int
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, literal_0539, http.NoBody)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code == http.StatusServiceUnavailable {
			drops++
		}
	}

	assert.True(t, drops >= 80, fmt.Sprintf("expected to be greater than 80, but got %d", drops))
}

const literal_6485 = "unit-test"

const literal_1932 = "X-Test"

const literal_0539 = "http://localhost"
